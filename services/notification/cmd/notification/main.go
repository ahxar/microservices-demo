package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/safar/microservices-demo/proto/notification/v1"
	"github.com/safar/microservices-demo/services/notification/internal/config"
	"github.com/safar/microservices-demo/services/notification/internal/repository"
	"github.com/safar/microservices-demo/services/notification/internal/server"
	"github.com/safar/microservices-demo/services/notification/internal/service"
	"github.com/safar/microservices-demo/services/notification/internal/smtp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository
	repo, err := repository.NewNotificationRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize SMTP client
	smtpClient := smtp.NewSMTPClient(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom)

	// Initialize service
	notificationService := service.NewNotificationService(repo, smtpClient)

	// Initialize gRPC server
	grpcServer := server.NewGRPCServer(notificationService)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create and configure gRPC server
	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, grpcServer)

	// Register reflection service for debugging
	reflection.Register(s)

	log.Printf("Notification Service starting on port %s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
