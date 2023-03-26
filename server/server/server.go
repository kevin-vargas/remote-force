package server

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	pb "remote-force/pb/V1"
	"remote-force/server/jwt"
	"remote-force/server/oauthdevice"
	"remote-force/server/store"

	"github.com/google/uuid"
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
func (s *server) Execute(ctx context.Context, cmdReq *pb.CommandRequest) (*pb.CommandResponse, error) {
	if _, ok := s.cmds[cmdReq.Name]; !ok {
		return &pb.CommandResponse{
			Type:   pb.Result_WARNING,
			Output: []byte(fmt.Sprintf("invalid command %s", cmdReq.Name)),
		}, nil
	}
	cmd := exec.CommandContext(ctx, cmdReq.Name, cmdReq.Args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return &pb.CommandResponse{
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
