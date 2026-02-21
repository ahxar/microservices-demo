package server

import (
	"context"
	"fmt"

	commonv1 "github.com/safar/microservices-demo/proto/common/v1"
	pb "github.com/safar/microservices-demo/proto/notification/v1"
	"github.com/safar/microservices-demo/services/notification/internal/service"
	"github.com/safar/microservices-demo/services/notification/internal/templates"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedNotificationServiceServer
	notificationService *service.NotificationService
}

func NewGRPCServer(notificationService *service.NotificationService) *GRPCServer {
	return &GRPCServer{
		notificationService: notificationService,
	}
}

func (s *GRPCServer) SendOrderConfirmation(ctx context.Context, req *pb.SendOrderConfirmationRequest) (*commonv1.Empty, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	// Convert proto items to template items
	var items []templates.OrderItem
	for _, item := range req.Items {
		items = append(items, templates.OrderItem{
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   formatMoney(item.UnitPrice),
			TotalPrice:  formatMoney(item.TotalPrice),
		})
	}

	// Prepare template data
	data := templates.OrderConfirmationData{
		FirstName: req.FirstName,
		OrderID:   req.OrderId,
		OrderDate: "Today", // Could be passed in request
		Items:     items,
		Subtotal:  formatMoney(req.Subtotal),
		Shipping:  formatMoney(req.Shipping),
		Tax:       "$0.00", // Not in request
		Total:     formatMoney(req.Total),
		ShippingAddress: templates.Address{
			Street:  req.ShippingAddress.Street,
			City:    req.ShippingAddress.City,
			State:   req.ShippingAddress.State,
			ZipCode: req.ShippingAddress.ZipCode,
			Country: req.ShippingAddress.Country,
		},
	}

	err := s.notificationService.SendOrderConfirmation(ctx, data, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send order confirmation: %v", err)
	}

	return &commonv1.Empty{}, nil
}

func (s *GRPCServer) SendShippingUpdate(ctx context.Context, req *pb.SendShippingUpdateRequest) (*commonv1.Empty, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	// Prepare template data
	data := templates.ShippingUpdateData{
		FirstName:      req.FirstName,
		OrderID:        req.OrderId,
		TrackingNumber: req.TrackingNumber,
		Carrier:        req.Carrier,
		Status:         req.Status,
		EstimatedDate:  "", // Could be calculated or passed
	}

	err := s.notificationService.SendShippingUpdate(ctx, data, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send shipping update: %v", err)
	}

	return &commonv1.Empty{}, nil
}

func (s *GRPCServer) SendWelcomeEmail(ctx context.Context, req *pb.SendWelcomeEmailRequest) (*commonv1.Empty, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Prepare template data
	data := templates.WelcomeData{
		FirstName: req.FirstName,
		Email:     req.Email,
	}

	err := s.notificationService.SendWelcomeEmail(ctx, data, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send welcome email: %v", err)
	}

	return &commonv1.Empty{}, nil
}

func (s *GRPCServer) SendPasswordReset(ctx context.Context, req *pb.SendPasswordResetRequest) (*commonv1.Empty, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.ResetToken == "" {
		return nil, status.Error(codes.InvalidArgument, "reset token is required")
	}

	// Prepare template data
	data := templates.PasswordResetData{
		FirstName:  "User", // Could be fetched from user service
		ResetToken: req.ResetToken,
		ResetURL:   fmt.Sprintf("http://localhost:3000/reset-password?token=%s", req.ResetToken),
		ExpiresIn:  "1 hour",
	}

	err := s.notificationService.SendPasswordReset(ctx, data, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send password reset: %v", err)
	}

	return &commonv1.Empty{}, nil
}

func formatMoney(money *commonv1.Money) string {
	if money == nil {
		return "$0.00"
	}
	dollars := float64(money.AmountCents) / 100.0
	return fmt.Sprintf("$%.2f", dollars)
}
