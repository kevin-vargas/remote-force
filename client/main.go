package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"remote-force/client/jwt"
	"remote-force/client/jwt/local"
	pb "remote-force/pb/V1"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	addr = "localhost:8080"
)

func red(f string, args ...interface{}) {
	color.Red(f, args...)
}

func info(f string, args ...interface{}) {
	color.Magenta(f, args...)
}

var (
	survey_want_login = &survey.Confirm{
		Message: "Do you want login now",
		Help:    "Login is a must to continue",
	}
)

type state struct {
	client pb.RemoteClient
	store  jwt.Store
}

func (s *state) start(ctx context.Context) error {
	o, ok, err := s.store.Get()
	if err != nil {
		return fmt.Errorf("failed to read store token %v", err)
	}
	if !ok {
		return s.login(ctx)
	}
	resCommands, err := s.client.Commands(ctx, &pb.CommandsRequest{})
	if err != nil {
		return fmt.Errorf("failed to send commands request %v", err)
	}
	askUsrInput := &survey.Input{
		Message: "Execute: ",
		Suggest: func(toComplete string) []string {
			var res []string
			for _, cmd := range resCommands.Commands {
				cmd = strings.ToLower(cmd)
				toComplete = strings.ToLower(toComplete)
				if strings.HasPrefix(cmd, toComplete) {
					res = append(res, cmd)
				}
			}
			return res
		},
		Help: "Available commands: " + fmt.Sprint(resCommands.Commands),
	}
	authCTX := authContext(ctx, o)
	for {
		var input string
		survey.AskOne(askUsrInput, &input, survey.WithIcons(func(is *survey.IconSet) {
			is.Question.Text = "😈"
		}))
		inputs := strings.Split(input, " ")
		if len(inputs) > 0 {
			name, args := inputs[0], inputs[1:]
			res, err := s.client.Execute(authCTX, &pb.ExecuteRequest{
				Name: name,
				Args: args,
			})
			if err != nil {
				if e, ok := status.FromError(err); ok {
					switch e.Code() {
					case codes.Unauthenticated:
						red(e.Message())
						s.cleanUP(ctx)
						return s.login(ctx)
					}
				}
				return err
			}
			out := string(res.Output)
			if res.Type == pb.Result_INFO {
				info(out)
			} else {
				fmt.Println(out)
			}
		}

	}
}

func (s *state) login(ctx context.Context) error {
	info("There is no active session")
	wantLog := false
	survey.AskOne(survey_want_login, &wantLog)
	if !wantLog {
		red("You must be logged in to continue")
		return errors.New("user dont want login")
	}
	res, err := s.client.Login(ctx, &pb.LoginRequest{})
	if err != nil {
		return fmt.Errorf("failed on call login %v", err)
	}
	info("Please log in")
	info("\tURL:\t%s", res.Url)
	info("\tCode:\t%s", res.Code)
	err = s.store.Save(res.Jwt)
	if err != nil {
		return fmt.Errorf("failed on save jwt token %v", err)
	}
	return s.start(ctx)
}

func (s *state) cleanUP(ctx context.Context) error {
	return s.store.CleanUP()
}

func authContext(ctx context.Context, o string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", o)
}

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}
	defer conn.Close()
	c := pb.NewRemoteClient(conn)
	s := local.New("./token")
	init := state{
		client: c,
		store:  s,
	}
	err = init.start(ctx)
	if err != nil {
		panic(err)
	}
}
