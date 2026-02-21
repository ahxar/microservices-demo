package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/safar/microservices-demo/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(serviceURL string) (*UserClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return &UserClient{
		client: pb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return c.client.Register(ctx, req)
}

func (c *UserClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	return c.client.Login(ctx, req)
}

func (c *UserClient) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	return c.client.RefreshToken(ctx, req)
}

func (c *UserClient) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, req)
}

func (c *UserClient) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	return c.client.GetUser(ctx, req)
}

func (c *UserClient) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	return c.client.UpdateUser(ctx, req)
}

func (c *UserClient) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return c.client.ListUsers(ctx, req)
}

func (c *UserClient) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.UserAddress, error) {
	return c.client.AddAddress(ctx, req)
}

func (c *UserClient) ListAddresses(ctx context.Context, req *pb.ListAddressesRequest) (*pb.ListAddressesResponse, error) {
	return c.client.ListAddresses(ctx, req)
}

func (c *UserClient) AddToWishlist(ctx context.Context, req *pb.AddToWishlistRequest) error {
	_, err := c.client.AddToWishlist(ctx, req)
	return err
}

func (c *UserClient) GetWishlist(ctx context.Context, req *pb.GetWishlistRequest) (*pb.WishlistResponse, error) {
	return c.client.GetWishlist(ctx, req)
}
