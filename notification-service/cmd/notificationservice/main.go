package main

import (
	"log"
	"net"
	"net/http"

	"github.com/NeGat1FF/e-commerce/notification-service/internal/config"
	"github.com/NeGat1FF/e-commerce/notification-service/internal/email"
	"github.com/NeGat1FF/e-commerce/notification-service/internal/server"
	mail "github.com/NeGat1FF/e-commerce/notification-service/proto"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.LoadConfig()

	client := &http.Client{}

	sender := email.NewEmailSender(client, cfg.JWTSecret)

	err := sender.InitTemplates("./templates/")
	if err != nil {
		panic(err)
	}

	server := server.NewServer(sender)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	mail.RegisterMailServiceServer(s, server)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
