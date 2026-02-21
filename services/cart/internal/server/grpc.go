package server

import (
	"context"

	commonv1 "github.com/safar/microservices-demo/proto/common/v1"
	pb "github.com/safar/microservices-demo/proto/cart/v1"
	"github.com/safar/microservices-demo/services/cart/internal/repository"
	"github.com/safar/microservices-demo/services/cart/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedCartServiceServer
	cartService *service.CartService
}

func NewGRPCServer(cartService *service.CartService) *GRPCServer {
	return &GRPCServer{
		cartService: cartService,
	}
}

func (s *GRPCServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.Cart, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	cart, err := s.cartService.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	return convertCartToProto(cart), nil
}

func (s *GRPCServer) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.Cart, error) {
	if req.UserId == "" || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and product ID are required")
	}

	if req.UnitPrice == nil {
		return nil, status.Error(codes.InvalidArgument, "unit price is required")
	}

	cart, err := s.cartService.AddItem(
		ctx,
		req.UserId,
		req.ProductId,
		req.ProductName,
		req.Quantity,
		repository.Money{
			AmountCents: req.UnitPrice.AmountCents,
			Currency:    req.UnitPrice.Currency,
		},
		req.ImageUrl,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add item: %v", err)
	}

	return convertCartToProto(cart), nil
}

func (s *GRPCServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.Cart, error) {
	if req.UserId == "" || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and product ID are required")
	}

	cart, err := s.cartService.UpdateItem(ctx, req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update item: %v", err)
	}

	return convertCartToProto(cart), nil
}

func (s *GRPCServer) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.Cart, error) {
	if req.UserId == "" || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and product ID are required")
	}

	cart, err := s.cartService.RemoveItem(ctx, req.UserId, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove item: %v", err)
	}

	return convertCartToProto(cart), nil
}

func (s *GRPCServer) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*commonv1.Empty, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if err := s.cartService.ClearCart(ctx, req.UserId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clear cart: %v", err)
	}

	return &commonv1.Empty{}, nil
}

func convertCartToProto(cart *repository.Cart) *pb.Cart {
	var items []*pb.CartItem
	var totalCents int64
	var currency string

	for _, item := range cart.Items {
		items = append(items, &pb.CartItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice: &commonv1.Money{
				AmountCents: item.UnitPrice.AmountCents,
				Currency:    item.UnitPrice.Currency,
			},
			TotalPrice: &commonv1.Money{
				AmountCents: item.TotalPrice.AmountCents,
				Currency:    item.TotalPrice.Currency,
			},
			ImageUrl: item.ImageURL,
		})

		totalCents += item.TotalPrice.AmountCents
		if currency == "" {
			currency = item.UnitPrice.Currency
		}
	}

	if currency == "" {
		currency = "USD"
	}

	return &pb.Cart{
		UserId: cart.UserID,
		Items:  items,
		Total: &commonv1.Money{
			AmountCents: totalCents,
			Currency:    currency,
		},
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
