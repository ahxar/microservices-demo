package client

import (
	"context"
	"fmt"
	"time"

	catalogpb "github.com/safar/microservices-demo/proto/catalog/v1"
	cartpb "github.com/safar/microservices-demo/proto/cart/v1"
	paymentpb "github.com/safar/microservices-demo/proto/payment/v1"
	shippingpb "github.com/safar/microservices-demo/proto/shipping/v1"
	notificationpb "github.com/safar/microservices-demo/proto/notification/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClients struct {
	Catalog      catalogpb.CatalogServiceClient
	Cart         cartpb.CartServiceClient
	Payment      paymentpb.PaymentServiceClient
	Shipping     shippingpb.ShippingServiceClient
	Notification notificationpb.NotificationServiceClient
	conns        []*grpc.ClientConn
}

func NewServiceClients(catalogURL, cartURL, paymentURL, shippingURL, notificationURL string) (*ServiceClients, error) {
	clients := &ServiceClients{}

	// Connect to Catalog Service
	catalogConn, err := dialService(catalogURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to catalog service: %w", err)
	}
	clients.Catalog = catalogpb.NewCatalogServiceClient(catalogConn)
	clients.conns = append(clients.conns, catalogConn)

	// Connect to Cart Service
	cartConn, err := dialService(cartURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cart service: %w", err)
	}
	clients.Cart = cartpb.NewCartServiceClient(cartConn)
	clients.conns = append(clients.conns, cartConn)

	// Connect to Payment Service
	paymentConn, err := dialService(paymentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	clients.Payment = paymentpb.NewPaymentServiceClient(paymentConn)
	clients.conns = append(clients.conns, paymentConn)

	// Connect to Shipping Service
	shippingConn, err := dialService(shippingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping service: %w", err)
	}
	clients.Shipping = shippingpb.NewShippingServiceClient(shippingConn)
	clients.conns = append(clients.conns, shippingConn)

	// Connect to Notification Service
	notificationConn, err := dialService(notificationURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}
	clients.Notification = notificationpb.NewNotificationServiceClient(notificationConn)
	clients.conns = append(clients.conns, notificationConn)

	return clients, nil
}

func dialService(serviceURL string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *ServiceClients) Close() error {
	for _, conn := range c.conns {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
