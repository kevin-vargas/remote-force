package main

import (
	"fmt"
	"log"
	"net"
	pb "remote-force/pb/V1"
	"remote-force/server/config"
	"remote-force/server/jwt"
	"remote-force/server/oauthdevice"
	"remote-force/server/server"
	"remote-force/server/server/store/memory"
	"time"

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
	// TODO: review time duration
	m := jwt.New(cfg.SecretKey, time.Hour*5000)
	ss := memory.New()
	o := oauthdevice.New(cfg.Oauth)
	rs := server.New(cfg.AvailableCommands, m, ss, o)
	s := grpc.NewServer()
	pb.RegisterRemoteServer(s, rs)
	if err = s.Serve(lis); err != nil {
		fmt.Printf("error on serve %v\n", err)
	}
}
