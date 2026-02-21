package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/safar/microservices-demo/proto/cart/v1"
	"github.com/safar/microservices-demo/services/cart/internal/config"
	"github.com/safar/microservices-demo/services/cart/internal/repository"
	"github.com/safar/microservices-demo/services/cart/internal/server"
	"github.com/safar/microservices-demo/services/cart/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository
	repo, err := repository.NewCartRepository(cfg.RedisURL, cfg.CartTTLDays)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize service
	cartService := service.NewCartService(repo)

	// Initialize gRPC server
	grpcServer := server.NewGRPCServer(cartService)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create and configure gRPC server
	s := grpc.NewServer()
	pb.RegisterCartServiceServer(s, grpcServer)

	// Register reflection service for debugging
	reflection.Register(s)

	log.Printf("Cart Service starting on port %s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
