package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sajanjswl/chat-service/proto"

	"google.golang.org/grpc/testdata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	ctx := context.Background()

	log.Println("startin chat client on port 8080...!")

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Unable to establish client connection to localhost:8080", err)
	}
	//defer conn.Close()
	client := proto.NewChatServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	stream, err := client.Chat(ctx)
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()

			if err != nil {
				log.Fatalf("Shutting down chat client : %v", err)
			}
			fmt.Println("Message from server ===> ", in.Message)
			message := strings.TrimRight(in.GetMessage(), "\n")
			if message == "BYE" || message == "bye" {

				close(waitc)

			}
		}
	}()

	go func() {

		for {
			text, _ := reader.ReadString('\n')
			err := stream.Send(&proto.Request{Message: text})
			if err != nil {
				log.Error(err)
			}

		}
	}()

	//	stream.CloseSend()
	<-waitc

}
