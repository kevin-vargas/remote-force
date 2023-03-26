package main

import (
	"context"
	"fmt"
	"log"
	pb "remote-force/pb/V1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addr = "localhost:8080"
)

func main() {
	fmt.Println("Im the client")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pb.NewRemoteClient(conn)

	// doExec(c)
	doLogin(c)
}

func doExec(c pb.RemoteClient) {

	res, err := c.Execute(context.Background(), &pb.CommandRequest{
		Name: "ls",
	})

	if err != nil {
		log.Fatalf("Failed on call %v\n", err)
		return
	}

	fmt.Println(res.Type)
	fmt.Println(string(res.Output))
}

func doLogin(c pb.RemoteClient) {

	res, err := c.Login(context.Background(), &pb.LoginRequest{})

	if err != nil {
		log.Fatalf("Failed on call %v\n", err)
		return
	}

	fmt.Println(res.Code)
	fmt.Println(res.Url)
	fmt.Println(res.Jwt)
}
