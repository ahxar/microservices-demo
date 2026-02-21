package service

import (
	"context"
	"fmt"

	"github.com/safar/microservices-demo/services/cart/internal/repository"
)

type CartService struct {
	repo *repository.CartRepository
}

func NewCartService(repo *repository.CartRepository) *CartService {
	return &CartService{
		repo: repo,
	}
}

func (s *CartService) GetCart(ctx context.Context, userID string) (*repository.Cart, error) {
	cart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Calculate total
	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, userID, productID, productName string, quantity int32, unitPrice repository.Money, imageURL string) (*repository.Cart, error) {
	item := repository.CartItem{
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		ImageURL:    imageURL,
	}

	cart, err := s.repo.AddItem(ctx, userID, item)
	if err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return cart, nil
}

func (s *CartService) UpdateItem(ctx context.Context, userID, productID string, quantity int32) (*repository.Cart, error) {
	cart, err := s.repo.UpdateItem(ctx, userID, productID, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return cart, nil
}

func (s *CartService) RemoveItem(ctx context.Context, userID, productID string) (*repository.Cart, error) {
	cart, err := s.repo.RemoveItem(ctx, userID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove item: %w", err)
	}

	return cart, nil
}

func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	if err := s.repo.DeleteCart(ctx, userID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}
