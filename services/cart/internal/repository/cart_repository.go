package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cart struct {
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int32  `json:"quantity"`
	UnitPrice   Money  `json:"unit_price"`
	TotalPrice  Money  `json:"total_price"`
	ImageURL    string `json:"image_url"`
}

type Money struct {
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

type CartRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCartRepository(redisURL string, ttlDays int) (*CartRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &CartRepository{
		client: client,
		ttl:    time.Duration(ttlDays) * 24 * time.Hour,
	}, nil
}

func (r *CartRepository) Close() error {
	return r.client.Close()
}

func (r *CartRepository) GetCart(ctx context.Context, userID string) (*Cart, error) {
	key := fmt.Sprintf("cart:%s", userID)

	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Return empty cart if not found
		return &Cart{
			UserID:    userID,
			Items:     []CartItem{},
			UpdatedAt: time.Now(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	var cart Cart
	if err := json.Unmarshal([]byte(data), &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	return &cart, nil
}

func (r *CartRepository) SaveCart(ctx context.Context, cart *Cart) error {
	key := fmt.Sprintf("cart:%s", cart.UserID)
	cart.UpdatedAt = time.Now()

	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	return nil
}

func (r *CartRepository) DeleteCart(ctx context.Context, userID string) error {
	key := fmt.Sprintf("cart:%s", userID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	return nil
}

func (r *CartRepository) AddItem(ctx context.Context, userID string, item CartItem) (*Cart, error) {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists
	found := false
	for i, existing := range cart.Items {
		if existing.ProductID == item.ProductID {
			// Update quantity
			cart.Items[i].Quantity += item.Quantity
			cart.Items[i].TotalPrice = Money{
				AmountCents: cart.Items[i].UnitPrice.AmountCents * int64(cart.Items[i].Quantity),
				Currency:    cart.Items[i].UnitPrice.Currency,
			}
			found = true
			break
		}
	}

	if !found {
		// Add new item
		item.TotalPrice = Money{
			AmountCents: item.UnitPrice.AmountCents * int64(item.Quantity),
			Currency:    item.UnitPrice.Currency,
		}
		cart.Items = append(cart.Items, item)
	}

	if err := r.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepository) UpdateItem(ctx context.Context, userID, productID string, quantity int32) (*Cart, error) {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			if quantity <= 0 {
				// Remove item if quantity is 0 or less
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i].Quantity = quantity
				cart.Items[i].TotalPrice = Money{
					AmountCents: cart.Items[i].UnitPrice.AmountCents * int64(quantity),
					Currency:    cart.Items[i].UnitPrice.Currency,
				}
			}
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("item not found in cart")
	}

	if err := r.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepository) RemoveItem(ctx context.Context, userID, productID string) (*Cart, error) {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}

	if err := r.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}
