package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/safar/microservices-demo/proto/catalog/v1"
	commonv1 "github.com/safar/microservices-demo/proto/common/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CatalogClient struct {
	client pb.CatalogServiceClient
	conn   *grpc.ClientConn
}

func NewCatalogClient(serviceURL string) (*CatalogClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to catalog service: %w", err)
	}

	return &CatalogClient{
		client: pb.NewCatalogServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *CatalogClient) Close() error {
	return c.conn.Close()
}

func (c *CatalogClient) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	return c.client.ListProducts(ctx, req)
}

func (c *CatalogClient) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	return c.client.GetProduct(ctx, req)
}

func (c *CatalogClient) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	return c.client.SearchProducts(ctx, req)
}

func (c *CatalogClient) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	return c.client.CreateProduct(ctx, req)
}

func (c *CatalogClient) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	return c.client.UpdateProduct(ctx, req)
}

func (c *CatalogClient) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) error {
	_, err := c.client.DeleteProduct(ctx, req)
	return err
}

func (c *CatalogClient) ListCategories(ctx context.Context) (*pb.ListCategoriesResponse, error) {
	return c.client.ListCategories(ctx, &commonv1.Empty{})
}
