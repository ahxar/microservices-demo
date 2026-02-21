package service

import (
	"context"
	"fmt"

	"github.com/safar/microservices-demo/services/catalog/internal/repository"
)

type CatalogService struct {
	repo *repository.CatalogRepository
}

func NewCatalogService(repo *repository.CatalogRepository) *CatalogService {
	return &CatalogService{
		repo: repo,
	}
}

// Category operations
func (s *CatalogService) ListCategories(ctx context.Context) ([]*repository.Category, error) {
	categories, err := s.repo.ListCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	return categories, nil
}

// Product operations
func (s *CatalogService) CreateProduct(ctx context.Context, name, slug, description string, priceCents int64, currency, categoryID string, imageURLs []string, stockQuantity int32) (*repository.Product, error) {
	product, err := s.repo.CreateProduct(name, slug, description, priceCents, currency, categoryID, imageURLs, stockQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}

func (s *CatalogService) GetProductByID(ctx context.Context, id string) (*repository.Product, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *CatalogService) GetProductBySlug(ctx context.Context, slug string) (*repository.Product, error) {
	product, err := s.repo.GetProductBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *CatalogService) ListProducts(ctx context.Context, page, pageSize int, categoryID string, activeOnly bool) ([]*repository.Product, int, error) {
	offset := (page - 1) * pageSize
	products, total, err := s.repo.ListProducts(pageSize, offset, categoryID, activeOnly)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	return products, total, nil
}

func (s *CatalogService) SearchProducts(ctx context.Context, query string, page, pageSize int, categoryID string) ([]*repository.Product, int, error) {
	offset := (page - 1) * pageSize
	products, total, err := s.repo.SearchProducts(query, pageSize, offset, categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}
	return products, total, nil
}

func (s *CatalogService) UpdateProduct(ctx context.Context, id, name, slug, description string, priceCents int64, currency, categoryID string, imageURLs []string, stockQuantity int32, isActive bool) (*repository.Product, error) {
	product, err := s.repo.UpdateProduct(id, name, slug, description, priceCents, currency, categoryID, imageURLs, stockQuantity, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return product, nil
}

func (s *CatalogService) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.DeleteProduct(id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// Inventory operations
func (s *CatalogService) CheckInventory(ctx context.Context, items map[string]int32) (bool, []string, error) {
	var unavailable []string

	for productID, quantity := range items {
		available, err := s.repo.CheckInventory(productID, quantity)
		if err != nil {
			return false, nil, fmt.Errorf("failed to check inventory for product %s: %w", productID, err)
		}
		if !available {
			unavailable = append(unavailable, productID)
		}
	}

	return len(unavailable) == 0, unavailable, nil
}

func (s *CatalogService) ReserveInventory(ctx context.Context, orderID string, items map[string]int32, expirationMinutes int32) (string, error) {
	// First check all items are available
	available, unavailable, err := s.CheckInventory(ctx, items)
	if err != nil {
		return "", err
	}
	if !available {
		return "", fmt.Errorf("insufficient inventory for products: %v", unavailable)
	}

	// Reserve each item
	var reservationID string
	for productID, quantity := range items {
		resID, err := s.repo.ReserveInventory(orderID, productID, quantity, expirationMinutes)
		if err != nil {
			return "", fmt.Errorf("failed to reserve inventory: %w", err)
		}
		if reservationID == "" {
			reservationID = resID
		}
	}

	return reservationID, nil
}
