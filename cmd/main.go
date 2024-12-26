package main

import (
	"log"
	"net"
	"whatsapp/db"
	routes "whatsapp/internal/routes"
	proto "whatsapp/proto/chatproto"
	"whatsapp/services"

	"google.golang.org/grpc"
)

func main() {
	db.GetDB()
	r := routes.SetUpRouter()
	r.Run(":8080")

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		proto.RegisterChatServiceServer(grpcServer, &services.ChatServiceServer{DB: db.DB})

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

}
