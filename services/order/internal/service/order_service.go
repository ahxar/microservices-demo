package service

import (
	"context"
	"database/sql"
	"fmt"

	commonpb "github.com/safar/microservices-demo/proto/common/v1"
	catalogpb "github.com/safar/microservices-demo/proto/catalog/v1"
	cartpb "github.com/safar/microservices-demo/proto/cart/v1"
	notificationpb "github.com/safar/microservices-demo/proto/notification/v1"
	paymentpb "github.com/safar/microservices-demo/proto/payment/v1"
	shippingpb "github.com/safar/microservices-demo/proto/shipping/v1"
	"github.com/safar/microservices-demo/services/order/internal/client"
	"github.com/safar/microservices-demo/services/order/internal/repository"
)

type OrderService struct {
	repo    *repository.OrderRepository
	clients *client.ServiceClients
}

func NewOrderService(repo *repository.OrderRepository, clients *client.ServiceClients) *OrderService {
	return &OrderService{
		repo:    repo,
		clients: clients,
	}
}

// CreateOrder orchestrates the entire checkout flow
func (s *OrderService) CreateOrder(ctx context.Context, userID string, shippingAddress *commonpb.Address, paymentMethodID string) (*repository.Order, []repository.OrderItem, error) {
	// Step 1: Get cart from Cart Service
	cart, err := s.clients.Cart.GetCart(ctx, &cartpb.GetCartRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(cart.Items) == 0 {
		return nil, nil, fmt.Errorf("cart is empty")
	}

	// Step 2: Validate products exist and check inventory
	var inventoryItems []*catalogpb.InventoryItem
	for _, item := range cart.Items {
		inventoryItems = append(inventoryItems, &catalogpb.InventoryItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	inventoryCheck, err := s.clients.Catalog.CheckInventory(ctx, &catalogpb.CheckInventoryRequest{
		Items: inventoryItems,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check inventory: %w", err)
	}

	if !inventoryCheck.Available {
		return nil, nil, fmt.Errorf("insufficient inventory for products: %v", inventoryCheck.UnavailableProductIds)
	}

	// Step 3: Reserve inventory
	reserveResp, err := s.clients.Catalog.ReserveInventory(ctx, &catalogpb.ReserveInventoryRequest{
		OrderId:           "", // Will be filled after order creation
		Items:             inventoryItems,
		ExpirationMinutes: 15,
	})
	if err != nil || !reserveResp.Success {
		return nil, nil, fmt.Errorf("failed to reserve inventory: %w", err)
	}

	// Step 4: Calculate shipping
	shippingQuote, err := s.clients.Shipping.GetQuote(ctx, &shippingpb.GetQuoteRequest{
		To:          shippingAddress,
		WeightGrams: 1000, // Mock weight
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get shipping quote: %w", err)
	}

	var shippingCost int64
	if len(shippingQuote.Quotes) > 0 {
		shippingCost = shippingQuote.Quotes[0].Cost.AmountCents
	}

	// Step 5: Calculate totals
	subtotal := cart.Total.AmountCents
	tax := int64(0) // Mock tax calculation
	total := subtotal + shippingCost + tax

	// Step 6: Process payment
	chargeResp, err := s.clients.Payment.Charge(ctx, &paymentpb.ChargeRequest{
		OrderId:         "",  // Will be filled after order creation
		UserId:          userID,
		PaymentMethodId: paymentMethodID,
		Amount: &commonpb.Money{
			AmountCents: total,
			Currency:    cart.Total.Currency,
		},
		IdempotencyKey: fmt.Sprintf("order-%s-%d", userID, subtotal),
	})
	if err != nil || !chargeResp.Success {
		return nil, nil, fmt.Errorf("payment failed: %v", chargeResp.ErrorMessage)
	}

	// Step 7: Create order record
	order := &repository.Order{
		UserID:          userID,
		Status:          "confirmed",
		SubtotalCents:   subtotal,
		ShippingCents:   shippingCost,
		TaxCents:        tax,
		TotalCents:      total,
		Currency:        cart.Total.Currency,
		ShippingStreet:  shippingAddress.Street,
		ShippingCity:    shippingAddress.City,
		ShippingState:   shippingAddress.State,
		ShippingZip:     shippingAddress.ZipCode,
		ShippingCountry: shippingAddress.Country,
		PaymentMethodID: sql.NullString{String: paymentMethodID, Valid: true},
		TransactionID:   sql.NullString{String: chargeResp.Transaction.Id, Valid: true},
	}

	var orderItems []repository.OrderItem
	for _, item := range cart.Items {
		orderItems = append(orderItems, repository.OrderItem{
			ProductID:       item.ProductId,
			ProductName:     item.ProductName,
			Quantity:        item.Quantity,
			UnitPriceCents:  item.UnitPrice.AmountCents,
			TotalPriceCents: item.TotalPrice.AmountCents,
		})
	}

	createdOrder, err := s.repo.CreateOrder(order, orderItems)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Step 8: Clear cart
	_, err = s.clients.Cart.ClearCart(ctx, &cartpb.ClearCartRequest{
		UserId: userID,
	})
	if err != nil {
		// Log error but don't fail order
		fmt.Printf("Warning: failed to clear cart: %v\n", err)
	}

	// Step 9: Send confirmation email (async - don't wait)
	go func() {
		var notifItems []*notificationpb.OrderItem
		for _, item := range orderItems {
			notifItems = append(notifItems, &notificationpb.OrderItem{
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				UnitPrice: &commonpb.Money{
					AmountCents: item.UnitPriceCents,
					Currency:    createdOrder.Currency,
				},
				TotalPrice: &commonpb.Money{
					AmountCents: item.TotalPriceCents,
					Currency:    createdOrder.Currency,
				},
			})
		}

		s.clients.Notification.SendOrderConfirmation(context.Background(), &notificationpb.SendOrderConfirmationRequest{
			Email:     "user@example.com", // Would come from user service
			FirstName: "User",
			OrderId:   createdOrder.ID,
			Items:     notifItems,
			Subtotal: &commonpb.Money{
				AmountCents: createdOrder.SubtotalCents,
				Currency:    createdOrder.Currency,
			},
			Shipping: &commonpb.Money{
				AmountCents: createdOrder.ShippingCents,
				Currency:    createdOrder.Currency,
			},
			Total: &commonpb.Money{
				AmountCents: createdOrder.TotalCents,
				Currency:    createdOrder.Currency,
			},
			ShippingAddress: shippingAddress,
		})
	}()

	return createdOrder, orderItems, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID, userID string) (*repository.Order, []repository.OrderItem, []repository.OrderStatusHistory, error) {
	order, items, history, err := s.repo.GetOrder(orderID, userID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, items, history, nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID string, page, pageSize int, statusFilter string) ([]*repository.Order, int, error) {
	offset := (page - 1) * pageSize
	orders, total, err := s.repo.ListOrders(userID, pageSize, offset, statusFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, total, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID, status, notes string) error {
	if err := s.repo.UpdateOrderStatus(orderID, status, notes); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID, reason string) error {
	// Get order to check status
	order, _, _, err := s.repo.GetOrder(orderID, userID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status == "shipped" || order.Status == "delivered" {
		return fmt.Errorf("cannot cancel order that has been shipped")
	}

	// Update status to cancelled
	if err := s.repo.UpdateOrderStatus(orderID, "cancelled", reason); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Trigger refund if payment was processed
	if order.TransactionID.Valid {
		_, err := s.clients.Payment.Refund(ctx, &paymentpb.RefundRequest{
			TransactionId: order.TransactionID.String,
			Amount: &commonpb.Money{
				AmountCents: order.TotalCents,
				Currency:    order.Currency,
			},
			Reason: reason,
		})
		if err != nil {
			// Log error but don't fail cancellation
			fmt.Printf("Warning: failed to process refund: %v\n", err)
		}
	}

	return nil
}
