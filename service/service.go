package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sajanjswl/chat-service/proto"
	log "github.com/sirupsen/logrus"
)

type chatServiceServer struct {
}

func NewChatServiceServer() proto.ChatServiceServer {
	return &chatServiceServer{}
}

func (chat *chatServiceServer) Chat(stream proto.ChatService_ChatServer) error {
	reader := bufio.NewReader(os.Stdin)
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()

			if err != nil {
				log.Fatalf("Shutting down chat server: %v", err)
			}
			fmt.Println("Message from client ====> ", in.Message)
			message := strings.TrimRight(in.GetMessage(), "\n")
			if message == "BYE" || message == "bye" {

				err := stream.Send(&proto.Response{Message: "bye"})
				if err != nil {
					log.Error(err)
				}
				close(waitc)

			}
		}
	}()

	go func() {

		for {
			text, _ := reader.ReadString('\n')
			err := stream.Send(&proto.Response{Message: text})
			if err != nil {
				log.Error(err)
			}

		}
	}()

	<-waitc
	return nil
}
