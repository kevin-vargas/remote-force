package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	pb "remote-force/pb/V1"
	"remote-force/server/interceptor"
	"remote-force/server/jwt"
	"remote-force/server/oauthdevice"
	"remote-force/server/store"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(cmds []string, m *jwt.Manager, s store.Token, o *oauthdevice.Config) pb.RemoteServer {
	serverCmds := make(map[string]any)
	for _, e := range cmds {
		serverCmds[e] = true
	}
	return &server{
		oauth:      o,
		tokenStore: s,
		cmds:       serverCmds,
		jwtManager: m,
	}
}

type server struct {
	tokenStore store.Token
	oauth      *oauthdevice.Config
	jwtManager *jwt.Manager
	cmds       map[string]any
	pb.RemoteServer
}

// TODO: improve error handler
// TODO: add log invalid command and execute command
func (s *server) Execute(ctx context.Context, cmdReq *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	usr, ok := interceptor.ContextUser(ctx)
	if !ok {
		return nil, status.Errorf(codes.Aborted, "invalid context on execute no user")
	}
	if _, ok := s.cmds[cmdReq.Name]; !ok {
		log.Printf("User: %s, try to invoke invalid command: %s With Args: %+v\n", usr.Email, cmdReq.Name, cmdReq.Args)
		return &pb.ExecuteResponse{
			Type:   pb.Result_WARNING,
			Output: []byte(fmt.Sprintf("invalid command %s\n", cmdReq.Name)),
		}, nil
	}

	cmd := exec.CommandContext(ctx, cmdReq.Name, cmdReq.Args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	log.Printf("User: %s, Execute: %s With Args: %+v\n", usr.Email, cmdReq.Name, cmdReq.Args)
	return &pb.ExecuteResponse{
		Type:   pb.Result_INFO,
		Output: out.Bytes(),
	}, nil
}

// TODO: validate before no token have client
func (s *server) Login(ctx context.Context, _ *pb.LoginRequest) (*pb.LoginResponse, error) {
	// TODO: user an interface to create ids
	id := uuid.New()
	j, err := s.jwtManager.Generate(id.String())
	if err != nil {
		return nil, err
	}
	o, err := s.oauth.AuthDevice(ctx)
	if err != nil {
		return nil, err
	}
	go s.PollToken(ctx, id.String(), o)
	return &pb.LoginResponse{
		Url:  o.VerificationURI,
		Code: o.UserCode,
		Jwt:  j,
	}, nil
}

func (s *server) Commands(ctx context.Context, _ *pb.CommandsRequest) (*pb.CommandsResponse, error) {
	var cmds []string
	for k := range s.cmds {
		cmds = append(cmds, k)
	}
	return &pb.CommandsResponse{
		Commands: cmds,
	}, nil
}

func (s *server) PollToken(ctx context.Context, id string, od *oauthdevice.DeviceAuth) {
	t, err := s.oauth.Poll(ctx, od)
	// TODO: log no print
	if err != nil {
		fmt.Println("on pull")
		fmt.Println(err.Error())
		return
	}
	if err = s.tokenStore.Save(id, t); err != nil {
		fmt.Println("on save")
		fmt.Println(err.Error())
		return
	}
}
