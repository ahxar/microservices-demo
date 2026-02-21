package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/safar/microservices-demo/proto/order/v1"
	"github.com/safar/microservices-demo/services/order/internal/client"
	"github.com/safar/microservices-demo/services/order/internal/config"
	"github.com/safar/microservices-demo/services/order/internal/repository"
	"github.com/safar/microservices-demo/services/order/internal/server"
	"github.com/safar/microservices-demo/services/order/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository
	repo, err := repository.NewOrderRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize service clients
	clients, err := client.NewServiceClients(
		cfg.CatalogServiceURL,
		cfg.CartServiceURL,
		cfg.PaymentServiceURL,
		cfg.ShippingServiceURL,
		cfg.NotificationServiceURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize service clients: %v", err)
	}
	defer clients.Close()

	// Initialize service
	orderService := service.NewOrderService(repo, clients)

	// Initialize gRPC server
	grpcServer := server.NewGRPCServer(orderService)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create and configure gRPC server
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, grpcServer)

	// Register reflection service for debugging
	reflection.Register(s)

	log.Printf("Order Service starting on port %s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
