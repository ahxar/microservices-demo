package server

import (
	"context"
	"math"

	commonv1 "github.com/safar/microservices-demo/proto/common/v1"
	pb "github.com/safar/microservices-demo/proto/order/v1"
	"github.com/safar/microservices-demo/services/order/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedOrderServiceServer
	orderService *service.OrderService
}

func NewGRPCServer(orderService *service.OrderService) *GRPCServer {
	return &GRPCServer{
		orderService: orderService,
	}
}

func (s *GRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	if req.UserId == "" || req.ShippingAddress == nil || req.PaymentMethodId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID, shipping address, and payment method ID are required")
	}

	order, items, err := s.orderService.CreateOrder(ctx, req.UserId, req.ShippingAddress, req.PaymentMethodId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Convert to proto
	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			Id:          item.ID,
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice: &commonv1.Money{
				AmountCents: item.UnitPriceCents,
				Currency:    order.Currency,
			},
			TotalPrice: &commonv1.Money{
				AmountCents: item.TotalPriceCents,
				Currency:    order.Currency,
			},
		})
	}

	pbOrder := &pb.Order{
		Id:     order.ID,
		UserId: order.UserID,
		Status: getOrderStatus(order.Status),
		Items:  pbItems,
		Subtotal: &commonv1.Money{
			AmountCents: order.SubtotalCents,
			Currency:    order.Currency,
		},
		Shipping: &commonv1.Money{
			AmountCents: order.ShippingCents,
			Currency:    order.Currency,
		},
		Tax: &commonv1.Money{
			AmountCents: order.TaxCents,
			Currency:    order.Currency,
		},
		Total: &commonv1.Money{
			AmountCents: order.TotalCents,
			Currency:    order.Currency,
		},
		ShippingAddress: &commonv1.Address{
			Street:  order.ShippingStreet,
			City:    order.ShippingCity,
			State:   order.ShippingState,
			ZipCode: order.ShippingZip,
			Country: order.ShippingCountry,
		},
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if order.PaymentMethodID.Valid {
		pbOrder.PaymentMethodId = order.PaymentMethodID.String
	}
	if order.TransactionID.Valid {
		pbOrder.TransactionId = order.TransactionID.String
	}
	if order.TrackingNumber.Valid {
		pbOrder.TrackingNumber = order.TrackingNumber.String
	}

	return pbOrder, nil
}

func (s *GRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID and user ID are required")
	}

	order, items, history, err := s.orderService.GetOrder(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	// Convert items
	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			Id:          item.ID,
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice: &commonv1.Money{
				AmountCents: item.UnitPriceCents,
				Currency:    order.Currency,
			},
			TotalPrice: &commonv1.Money{
				AmountCents: item.TotalPriceCents,
				Currency:    order.Currency,
			},
		})
	}

	// Convert history
	var pbHistory []*pb.OrderStatusHistory
	for _, h := range history {
		pbHistory = append(pbHistory, &pb.OrderStatusHistory{
			Id:        h.ID,
			Status:    getOrderStatus(h.Status),
			Notes:     h.Notes,
			CreatedAt: h.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	pbOrder := &pb.Order{
		Id:     order.ID,
		UserId: order.UserID,
		Status: getOrderStatus(order.Status),
		Items:  pbItems,
		Subtotal: &commonv1.Money{
			AmountCents: order.SubtotalCents,
			Currency:    order.Currency,
		},
		Shipping: &commonv1.Money{
			AmountCents: order.ShippingCents,
			Currency:    order.Currency,
		},
		Tax: &commonv1.Money{
			AmountCents: order.TaxCents,
			Currency:    order.Currency,
		},
		Total: &commonv1.Money{
			AmountCents: order.TotalCents,
			Currency:    order.Currency,
		},
		ShippingAddress: &commonv1.Address{
			Street:  order.ShippingStreet,
			City:    order.ShippingCity,
			State:   order.ShippingState,
			ZipCode: order.ShippingZip,
			Country: order.ShippingCountry,
		},
		History:   pbHistory,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if order.PaymentMethodID.Valid {
		pbOrder.PaymentMethodId = order.PaymentMethodID.String
	}
	if order.TransactionID.Valid {
		pbOrder.TransactionId = order.TransactionID.String
	}
	if order.TrackingNumber.Valid {
		pbOrder.TrackingNumber = order.TrackingNumber.String
	}

	return pbOrder, nil
}

func (s *GRPCServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	page := int(req.Pagination.Page)
	pageSize := int(req.Pagination.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	statusFilter := ""
	if req.StatusFilter != pb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		statusFilter = getOrderStatusString(req.StatusFilter)
	}

	orders, total, err := s.orderService.ListOrders(ctx, req.UserId, page, pageSize, statusFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrder := &pb.Order{
			Id:     order.ID,
			UserId: order.UserID,
			Status: getOrderStatus(order.Status),
			Subtotal: &commonv1.Money{
				AmountCents: order.SubtotalCents,
				Currency:    order.Currency,
			},
			Shipping: &commonv1.Money{
				AmountCents: order.ShippingCents,
				Currency:    order.Currency,
			},
			Tax: &commonv1.Money{
				AmountCents: order.TaxCents,
				Currency:    order.Currency,
			},
			Total: &commonv1.Money{
				AmountCents: order.TotalCents,
				Currency:    order.Currency,
			},
			CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		pbOrders = append(pbOrders, pbOrder)
	}

	totalPages := int32(math.Ceil(float64(total) / float64(pageSize)))

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
		Pagination: &commonv1.PaginationResponse{
			Page:       int32(page),
			PageSize:   int32(pageSize),
			TotalPages: totalPages,
			TotalCount: int64(total),
		},
	}, nil
}

func (s *GRPCServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	statusStr := getOrderStatusString(req.Status)
	if err := s.orderService.UpdateOrderStatus(ctx, req.Id, statusStr, req.Notes); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	// Return updated order (simplified - would need user ID)
	return &pb.Order{
		Id:     req.Id,
		Status: req.Status,
	}, nil
}

func (s *GRPCServer) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.Order, error) {
	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID and user ID are required")
	}

	if err := s.orderService.CancelOrder(ctx, req.Id, req.UserId, req.Reason); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	return &pb.Order{
		Id:     req.Id,
		Status: pb.OrderStatus_ORDER_STATUS_CANCELLED,
	}, nil
}

func getOrderStatus(status string) pb.OrderStatus {
	switch status {
	case "pending":
		return pb.OrderStatus_ORDER_STATUS_PENDING
	case "confirmed":
		return pb.OrderStatus_ORDER_STATUS_CONFIRMED
	case "processing":
		return pb.OrderStatus_ORDER_STATUS_PROCESSING
	case "shipped":
		return pb.OrderStatus_ORDER_STATUS_SHIPPED
	case "delivered":
		return pb.OrderStatus_ORDER_STATUS_DELIVERED
	case "cancelled":
		return pb.OrderStatus_ORDER_STATUS_CANCELLED
	case "refunded":
		return pb.OrderStatus_ORDER_STATUS_REFUNDED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func getOrderStatusString(status pb.OrderStatus) string {
	switch status {
	case pb.OrderStatus_ORDER_STATUS_PENDING:
		return "pending"
	case pb.OrderStatus_ORDER_STATUS_CONFIRMED:
		return "confirmed"
	case pb.OrderStatus_ORDER_STATUS_PROCESSING:
		return "processing"
	case pb.OrderStatus_ORDER_STATUS_SHIPPED:
		return "shipped"
	case pb.OrderStatus_ORDER_STATUS_DELIVERED:
		return "delivered"
	case pb.OrderStatus_ORDER_STATUS_CANCELLED:
		return "cancelled"
	case pb.OrderStatus_ORDER_STATUS_REFUNDED:
		return "refunded"
	default:
		return "pending"
	}
}
