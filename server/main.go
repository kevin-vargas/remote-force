package main

import (
	"fmt"
	"log"
	"net"
	pb "remote-force/pb/V1"
	"remote-force/server/config"
	"remote-force/server/entity"
	"remote-force/server/interceptor"
	"remote-force/server/interceptor/user-provider/github"
	"remote-force/server/jwt"
	"remote-force/server/oauthdevice"
	"remote-force/server/server"
	"remote-force/server/store/memory"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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
	ts := memory.New[*oauth2.Token]()
	o := oauthdevice.New(cfg.Oauth)
	rs := server.New(cfg.AvailableCommands, m, ts, o)
	ai := interceptor.Authentication(m)
	up := github.New(ts)
	us := memory.New[entity.User]()
	ui := interceptor.UserInfo(up, us)
	i := grpc.ChainUnaryInterceptor(ai, ui)
	s := grpc.NewServer(i)
	pb.RegisterRemoteServer(s, rs)
	if err = s.Serve(lis); err != nil {
		fmt.Printf("error on serve %v\n", err)
	}
}
