package service

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	commonpb "github.com/safar/microservices-demo/proto/common/v1"
	"github.com/safar/microservices-demo/services/shipping/internal/repository"
)

type ShippingService struct {
	repo *repository.ShippingRepository
}

type ShippingQuote struct {
	Carrier       string
	Service       string
	CostCents     int64
	Currency      string
	EstimatedDays int32
}

// Shipment is an exported alias for repository.Shipment
type Shipment = repository.Shipment

func NewShippingService(repo *repository.ShippingRepository) *ShippingService {
	rand.Seed(time.Now().UnixNano())
	return &ShippingService{repo: repo}
}

// GetQuote generates mock shipping quotes based on address and weight
func (s *ShippingService) GetQuote(ctx context.Context, from, to *commonpb.Address, weightGrams int32) ([]ShippingQuote, error) {
	if to == nil {
		return nil, fmt.Errorf("destination address is required")
	}

	if weightGrams <= 0 {
		return nil, fmt.Errorf("weight must be positive")
	}

	// Mock shipping quotes (in real system, would call carrier APIs)
	quotes := []ShippingQuote{
		{
			Carrier:       "USPS",
			Service:       "Priority Mail",
			CostCents:     calculateShippingCost(weightGrams, 0.15),
			Currency:      "USD",
			EstimatedDays: 3,
		},
		{
			Carrier:       "USPS",
			Service:       "First Class",
			CostCents:     calculateShippingCost(weightGrams, 0.10),
			Currency:      "USD",
			EstimatedDays: 5,
		},
		{
			Carrier:       "FedEx",
			Service:       "Ground",
			CostCents:     calculateShippingCost(weightGrams, 0.18),
			Currency:      "USD",
			EstimatedDays: 4,
		},
		{
			Carrier:       "FedEx",
			Service:       "2-Day",
			CostCents:     calculateShippingCost(weightGrams, 0.35),
			Currency:      "USD",
			EstimatedDays: 2,
		},
		{
			Carrier:       "UPS",
			Service:       "Ground",
			CostCents:     calculateShippingCost(weightGrams, 0.17),
			Currency:      "USD",
			EstimatedDays: 4,
		},
	}

	return quotes, nil
}

// CreateShipment creates a new shipment with tracking
func (s *ShippingService) CreateShipment(ctx context.Context, orderID string, from, to *commonpb.Address, weightGrams int32, carrier, service string) (*repository.Shipment, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	if from == nil || to == nil {
		return nil, fmt.Errorf("from and to addresses are required")
	}

	// Generate tracking number
	trackingNumber := generateTrackingNumber(carrier)

	// Get shipping quote for selected carrier/service
	quotes, err := s.GetQuote(ctx, from, to, weightGrams)
	if err != nil {
		return nil, err
	}

	var selectedQuote *ShippingQuote
	for _, quote := range quotes {
		if quote.Carrier == carrier && quote.Service == service {
			selectedQuote = &quote
			break
		}
	}

	if selectedQuote == nil {
		// Use first quote as fallback
		selectedQuote = &quotes[0]
	}

	shipment := &repository.Shipment{
		OrderID:           orderID,
		TrackingNumber:    trackingNumber,
		Carrier:           selectedQuote.Carrier,
		Service:           selectedQuote.Service,
		Status:            "pending",
		FromStreet:        from.Street,
		FromCity:          from.City,
		FromState:         from.State,
		FromZip:           from.ZipCode,
		FromCountry:       from.Country,
		ToStreet:          to.Street,
		ToCity:            to.City,
		ToState:           to.State,
		ToZip:             to.ZipCode,
		ToCountry:         to.Country,
		WeightGrams:       weightGrams,
		ShippingCostCents: selectedQuote.CostCents,
		Currency:          selectedQuote.Currency,
		EstimatedDays:     sql.NullInt32{Int32: selectedQuote.EstimatedDays, Valid: true},
	}

	createdShipment, err := s.repo.CreateShipment(shipment)
	if err != nil {
		return nil, err
	}

	// Add initial tracking event
	event := &repository.TrackingEvent{
		ShipmentID:  createdShipment.ID,
		Status:      "pending",
		Location:    from.City + ", " + from.State,
		Description: "Shipment created",
	}
	_ = s.repo.AddTrackingEvent(event)

	return createdShipment, nil
}

// GetShipmentByTracking retrieves shipment by tracking number
func (s *ShippingService) GetShipmentByTracking(ctx context.Context, trackingNumber string) (*repository.Shipment, []repository.TrackingEvent, error) {
	shipment, err := s.repo.GetShipmentByTracking(trackingNumber)
	if err != nil {
		return nil, nil, err
	}

	events, err := s.repo.GetTrackingEvents(shipment.ID)
	if err != nil {
		return nil, nil, err
	}

	return shipment, events, nil
}

// GetShipmentByID retrieves shipment by shipment ID
func (s *ShippingService) GetShipmentByID(ctx context.Context, shipmentID string) (*repository.Shipment, []repository.TrackingEvent, error) {
	shipment, err := s.repo.GetShipmentByID(shipmentID)
	if err != nil {
		return nil, nil, err
	}

	events, err := s.repo.GetTrackingEvents(shipment.ID)
	if err != nil {
		return nil, nil, err
	}

	return shipment, events, nil
}

// GetShipmentByOrderID retrieves shipment by order ID
func (s *ShippingService) GetShipmentByOrderID(ctx context.Context, orderID string) (*repository.Shipment, []repository.TrackingEvent, error) {
	shipment, err := s.repo.GetShipmentByOrderID(orderID)
	if err != nil {
		return nil, nil, err
	}

	events, err := s.repo.GetTrackingEvents(shipment.ID)
	if err != nil {
		return nil, nil, err
	}

	return shipment, events, nil
}

// UpdateShipmentStatus updates shipment status and adds tracking event
func (s *ShippingService) UpdateShipmentStatus(ctx context.Context, shipmentID, status, location, description string) error {
	if err := s.repo.UpdateShipmentStatus(shipmentID, status); err != nil {
		return err
	}

	event := &repository.TrackingEvent{
		ShipmentID:  shipmentID,
		Status:      status,
		Location:    location,
		Description: description,
	}

	return s.repo.AddTrackingEvent(event)
}

// Helper functions

func calculateShippingCost(weightGrams int32, ratePerGram float64) int64 {
	// Base cost + weight-based cost
	baseCost := 500 // $5.00 base
	weightCost := int64(float64(weightGrams) * ratePerGram)
	return int64(baseCost) + weightCost
}

func generateTrackingNumber(carrier string) string {
	// Generate mock tracking numbers in carrier-specific formats
	timestamp := time.Now().Unix()
	random := rand.Intn(999999)

	switch carrier {
	case "USPS":
		return fmt.Sprintf("9400%d%06d", timestamp%100000000, random)
	case "FedEx":
		return fmt.Sprintf("%d%06d", timestamp%1000000000, random)
	case "UPS":
		return fmt.Sprintf("1Z%d%06d", timestamp%10000000, random)
	default:
		return fmt.Sprintf("SHIP%d%06d", timestamp%100000000, random)
	}
}
