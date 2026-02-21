package server

import (
	"context"

	commonv1 "github.com/safar/microservices-demo/proto/common/v1"
	pb "github.com/safar/microservices-demo/proto/shipping/v1"
	"github.com/safar/microservices-demo/services/shipping/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedShippingServiceServer
	shippingService *service.ShippingService
}

func NewGRPCServer(shippingService *service.ShippingService) *GRPCServer {
	return &GRPCServer{
		shippingService: shippingService,
	}
}

func (s *GRPCServer) GetQuote(ctx context.Context, req *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	if req.To == nil {
		return nil, status.Error(codes.InvalidArgument, "destination address is required")
	}

	if req.WeightGrams <= 0 {
		return nil, status.Error(codes.InvalidArgument, "weight must be positive")
	}

	quotes, err := s.shippingService.GetQuote(ctx, req.From, req.To, req.WeightGrams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get quotes: %v", err)
	}

	var pbQuotes []*pb.ShippingQuote
	for _, quote := range quotes {
		pbQuotes = append(pbQuotes, &pb.ShippingQuote{
			Carrier: quote.Carrier,
			Service: quote.Service,
			Cost: &commonv1.Money{
				AmountCents: quote.CostCents,
				Currency:    quote.Currency,
			},
			EstimatedDays: quote.EstimatedDays,
		})
	}

	return &pb.GetQuoteResponse{
		Quotes: pbQuotes,
	}, nil
}

func (s *GRPCServer) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.Shipment, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	if req.From == nil || req.To == nil {
		return nil, status.Error(codes.InvalidArgument, "from and to addresses are required")
	}

	shipment, err := s.shippingService.CreateShipment(
		ctx,
		req.OrderId,
		req.From,
		req.To,
		req.WeightGrams,
		req.Carrier,
		req.Service,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create shipment: %v", err)
	}

	return &pb.Shipment{
		Id:             shipment.ID,
		OrderId:        shipment.OrderID,
		TrackingNumber: shipment.TrackingNumber,
		Carrier:        shipment.Carrier,
		Service:        shipment.Service,
		Status:         getShipmentStatus(shipment.Status),
		CreatedAt:      shipment.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GRPCServer) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.Shipment, error) {
	var shipment *service.Shipment
	var err error

	switch id := req.Identifier.(type) {
	case *pb.GetShipmentRequest_Id:
		shipment, _, err = s.shippingService.GetShipmentByID(ctx, id.Id)
	case *pb.GetShipmentRequest_OrderId:
		shipment, _, err = s.shippingService.GetShipmentByOrderID(ctx, id.OrderId)
	case *pb.GetShipmentRequest_TrackingNumber:
		shipment, _, err = s.shippingService.GetShipmentByTracking(ctx, id.TrackingNumber)
	default:
		return nil, status.Error(codes.InvalidArgument, "identifier is required")
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shipment not found: %v", err)
	}

	return &pb.Shipment{
		Id:             shipment.ID,
		OrderId:        shipment.OrderID,
		TrackingNumber: shipment.TrackingNumber,
		Carrier:        shipment.Carrier,
		Service:        shipment.Service,
		Status:         getShipmentStatus(shipment.Status),
		CreatedAt:      shipment.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GRPCServer) TrackShipment(ctx context.Context, req *pb.TrackShipmentRequest) (*pb.TrackShipmentResponse, error) {
	if req.TrackingNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "tracking number is required")
	}

	shipment, events, err := s.shippingService.GetShipmentByTracking(ctx, req.TrackingNumber)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shipment not found: %v", err)
	}

	pbShipment := &pb.Shipment{
		Id:             shipment.ID,
		OrderId:        shipment.OrderID,
		TrackingNumber: shipment.TrackingNumber,
		Carrier:        shipment.Carrier,
		Service:        shipment.Service,
		Status:         getShipmentStatus(shipment.Status),
		CreatedAt:      shipment.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	var pbEvents []*pb.TrackingEvent
	for _, event := range events {
		pbEvents = append(pbEvents, &pb.TrackingEvent{
			Id:          event.ID,
			Description: event.Description,
			Location:    event.Location,
			Timestamp:   event.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &pb.TrackShipmentResponse{
		Shipment: pbShipment,
		Events:   pbEvents,
	}, nil
}

func getShipmentStatus(status string) pb.ShipmentStatus {
	switch status {
	case "pending":
		return pb.ShipmentStatus_SHIPMENT_STATUS_PENDING
	case "label_created":
		return pb.ShipmentStatus_SHIPMENT_STATUS_LABEL_CREATED
	case "picked_up":
		return pb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP
	case "in_transit":
		return pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case "out_for_delivery":
		return pb.ShipmentStatus_SHIPMENT_STATUS_OUT_FOR_DELIVERY
	case "delivered":
		return pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	case "failed":
		return pb.ShipmentStatus_SHIPMENT_STATUS_FAILED
	default:
		return pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}
