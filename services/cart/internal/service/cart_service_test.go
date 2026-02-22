package service

import (
	"context"
	"errors"
	"testing"

	"github.com/safar/microservices-demo/services/cart/internal/repository"
)

type mockCartStore struct {
	addItemFn    func(ctx context.Context, userID string, item repository.CartItem) (*repository.Cart, error)
	deleteCartFn func(ctx context.Context, userID string) error
}

func (m *mockCartStore) GetCart(ctx context.Context, userID string) (*repository.Cart, error) {
	return &repository.Cart{UserID: userID}, nil
}

func (m *mockCartStore) AddItem(ctx context.Context, userID string, item repository.CartItem) (*repository.Cart, error) {
	if m.addItemFn != nil {
		return m.addItemFn(ctx, userID, item)
	}
	return &repository.Cart{UserID: userID, Items: []repository.CartItem{item}}, nil
}

func (m *mockCartStore) UpdateItem(ctx context.Context, userID, productID string, quantity int32) (*repository.Cart, error) {
	return &repository.Cart{UserID: userID}, nil
}

func (m *mockCartStore) RemoveItem(ctx context.Context, userID, productID string) (*repository.Cart, error) {
	return &repository.Cart{UserID: userID}, nil
}

func (m *mockCartStore) DeleteCart(ctx context.Context, userID string) error {
	if m.deleteCartFn != nil {
		return m.deleteCartFn(ctx, userID)
	}
	return nil
}

func TestAddItemBuildsRepositoryItem(t *testing.T) {
	var captured repository.CartItem
	store := &mockCartStore{
		addItemFn: func(_ context.Context, userID string, item repository.CartItem) (*repository.Cart, error) {
			if userID != "user-1" {
				t.Fatalf("unexpected user id: %s", userID)
			}
			captured = item
			return &repository.Cart{UserID: userID, Items: []repository.CartItem{item}}, nil
		},
	}

	svc := NewCartService(store)
	_, err := svc.AddItem(
		context.Background(),
		"user-1",
		"prod-1",
		"Widget",
		2,
		repository.Money{AmountCents: 1599, Currency: "USD"},
		"https://example.test/p.png",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if captured.ProductID != "prod-1" || captured.Quantity != 2 || captured.ProductName != "Widget" {
		t.Fatalf("unexpected captured item: %+v", captured)
	}
}

func TestClearCartWrapsDeleteError(t *testing.T) {
	store := &mockCartStore{
		deleteCartFn: func(_ context.Context, _ string) error {
			return errors.New("db down")
		},
	}

	svc := NewCartService(store)
	err := svc.ClearCart(context.Background(), "user-1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "failed to clear cart: db down" {
		t.Fatalf("unexpected error message: %v", err)
	}
}
