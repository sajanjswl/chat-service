package grpc

import (
	"context"

	"net"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/sajanjswl/chat-service/proto"

	"google.golang.org/grpc/testdata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func RunServer(ctx context.Context, chatAPI proto.ChatServiceServer, port string) {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {

	}
	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	// register service
	server := grpc.NewServer(grpc.Creds(creds))
	proto.RegisterChatServiceServer(server, chatAPI)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	server.Serve(listen)
}
