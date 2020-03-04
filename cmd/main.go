package main

import (
	"context"
	"flag"

	protocol "github.com/sajanjswl/chat-service/server/grpc"
	"github.com/sajanjswl/chat-service/service"
)

type Config struct {
	GRPCPort string
}

func main() {

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.Parse()

	ctx := context.Background()

	chatAPI := service.NewChatServiceServer()
	protocol.RunServer(ctx, chatAPI, cfg.GRPCPort)
}
