package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/safar/microservices-demo/proto/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	client pb.OrderServiceClient
	conn   *grpc.ClientConn
}

func NewOrderClient(serviceURL string) (*OrderClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	return &OrderClient{
		client: pb.NewOrderServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *OrderClient) Close() error {
	return c.conn.Close()
}

func (c *OrderClient) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	return c.client.CreateOrder(ctx, req)
}

func (c *OrderClient) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return c.client.GetOrder(ctx, req)
}

func (c *OrderClient) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	return c.client.ListOrders(ctx, req)
}

func (c *OrderClient) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.Order, error) {
	return c.client.CancelOrder(ctx, req)
}
