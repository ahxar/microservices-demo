package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/safar/microservices-demo/proto/cart/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CartClient struct {
	client pb.CartServiceClient
	conn   *grpc.ClientConn
}

func NewCartClient(serviceURL string) (*CartClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cart service: %w", err)
	}

	return &CartClient{
		client: pb.NewCartServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *CartClient) Close() error {
	return c.conn.Close()
}

func (c *CartClient) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.Cart, error) {
	return c.client.GetCart(ctx, req)
}

func (c *CartClient) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.Cart, error) {
	return c.client.AddItem(ctx, req)
}

func (c *CartClient) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.Cart, error) {
	return c.client.UpdateItem(ctx, req)
}

func (c *CartClient) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.Cart, error) {
	return c.client.RemoveItem(ctx, req)
}

func (c *CartClient) ClearCart(ctx context.Context, req *pb.ClearCartRequest) error {
	_, err := c.client.ClearCart(ctx, req)
	return err
}
