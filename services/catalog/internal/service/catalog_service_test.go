package service

import (
	"context"
	"errors"
	"testing"

	"github.com/safar/microservices-demo/services/catalog/internal/repository"
)

type mockCatalogRepository struct {
	listProductsFn     func(limit, offset int, categoryID string, activeOnly bool) error
	checkInventoryFn   func(productID string, quantity int32) (bool, error)
	reserveInventoryFn func(orderID, productID string, quantity, expirationMinutes int32) (string, error)
}

func (m *mockCatalogRepository) ListCategories() ([]*repository.Category, error) {
	return nil, nil
}

func (m *mockCatalogRepository) CreateProduct(name, slug, description string, priceCents int64, currency, categoryID string, imageURLs []string, stockQuantity int32) (*repository.Product, error) {
	return nil, nil
}

func (m *mockCatalogRepository) GetProductByID(id string) (*repository.Product, error) {
	return nil, nil
}

func (m *mockCatalogRepository) GetProductBySlug(slug string) (*repository.Product, error) {
	return nil, nil
}

func (m *mockCatalogRepository) ListProducts(limit, offset int, categoryID string, activeOnly bool) ([]*repository.Product, int, error) {
	if m.listProductsFn != nil {
		if err := m.listProductsFn(limit, offset, categoryID, activeOnly); err != nil {
			return nil, 0, err
		}
	}
	return []*repository.Product{}, 0, nil
}

func (m *mockCatalogRepository) SearchProducts(searchQuery string, limit, offset int, categoryID string) ([]*repository.Product, int, error) {
	return nil, 0, nil
}

func (m *mockCatalogRepository) UpdateProduct(id, name, slug, description string, priceCents int64, currency, categoryID string, imageURLs []string, stockQuantity int32, isActive bool) (*repository.Product, error) {
	return nil, nil
}

func (m *mockCatalogRepository) DeleteProduct(id string) error {
	return nil
}

func (m *mockCatalogRepository) CheckInventory(productID string, quantity int32) (bool, error) {
	if m.checkInventoryFn != nil {
		return m.checkInventoryFn(productID, quantity)
	}
	return true, nil
}

func (m *mockCatalogRepository) ReserveInventory(orderID, productID string, quantity, expirationMinutes int32) (string, error) {
	if m.reserveInventoryFn != nil {
		return m.reserveInventoryFn(orderID, productID, quantity, expirationMinutes)
	}
	return "res-default", nil
}

func TestListProductsCalculatesOffset(t *testing.T) {
	mockRepo := &mockCatalogRepository{
		listProductsFn: func(limit, offset int, categoryID string, activeOnly bool) error {
			if limit != 20 || offset != 40 {
				return errors.New("unexpected pagination values")
			}
			if categoryID != "cat-1" || !activeOnly {
				return errors.New("unexpected filter values")
			}
			return nil
		},
	}

	svc := NewCatalogService(mockRepo)
	_, _, err := svc.ListProducts(context.Background(), 3, 20, "cat-1", true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheckInventoryReturnsUnavailableItems(t *testing.T) {
	mockRepo := &mockCatalogRepository{
		checkInventoryFn: func(productID string, quantity int32) (bool, error) {
			if productID == "prod-2" {
				return false, nil
			}
			return true, nil
		},
	}

	svc := NewCatalogService(mockRepo)
	available, unavailable, err := svc.CheckInventory(context.Background(), map[string]int32{
		"prod-1": 1,
		"prod-2": 2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if available {
		t.Fatalf("expected inventory to be unavailable")
	}
	if len(unavailable) != 1 || unavailable[0] != "prod-2" {
		t.Fatalf("unexpected unavailable list: %#v", unavailable)
	}
}

func TestReserveInventoryUsesFirstReservationID(t *testing.T) {
	var calls int
	mockRepo := &mockCatalogRepository{
		checkInventoryFn: func(productID string, quantity int32) (bool, error) {
			return true, nil
		},
		reserveInventoryFn: func(orderID, productID string, quantity, expirationMinutes int32) (string, error) {
			calls++
			if calls == 1 {
				return "res-1", nil
			}
			return "res-2", nil
		},
	}

	svc := NewCatalogService(mockRepo)
	reservationID, err := svc.ReserveInventory(context.Background(), "order-1", map[string]int32{
		"prod-1": 1,
		"prod-2": 2,
	}, 15)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if reservationID != "res-1" {
		t.Fatalf("expected first reservation id, got %s", reservationID)
	}
}
