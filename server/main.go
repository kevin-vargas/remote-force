package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	pb "remote-force/pb/V1"
	"remote-force/server/config"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", cfg.ADDR)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rs := NewServer(cfg)
	pb.RegisterRemoteServer(s, rs)
	if err = s.Serve(lis); err != nil {
		fmt.Printf("error on serve %v\n", err)
	}
}

func NewServer(cfg config.Config) pb.RemoteServer {
	cmds := make(map[string]any)
	for _, e := range cfg.AvailableCommands {
		cmds[e] = true
	}
	return &server{
		cmds: cmds,
	}
}

type server struct {
	cmds map[string]any
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
