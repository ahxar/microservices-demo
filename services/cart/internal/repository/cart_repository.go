package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
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
	db *sql.DB
}

func NewCartRepository(databaseURL string) (*CartRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &CartRepository{db: db}, nil
}

func (r *CartRepository) Close() error {
	return r.db.Close()
}

func (r *CartRepository) GetCart(ctx context.Context, userID string) (*Cart, error) {
	cart := &Cart{
		UserID: userID,
		Items:  []CartItem{},
	}

	err := r.db.QueryRowContext(ctx, `
		SELECT updated_at
		FROM carts
		WHERE user_id = $1
	`, userID).Scan(&cart.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cart.UpdatedAt = time.Now().UTC()
			return cart, nil
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT product_id, product_name, quantity, unit_price_cents, currency, image_url
		FROM cart_items
		WHERE user_id = $1
		ORDER BY product_id ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := CartItem{}
		if err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.UnitPrice.AmountCents,
			&item.UnitPrice.Currency,
			&item.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		item.TotalPrice = Money{
			AmountCents: item.UnitPrice.AmountCents * int64(item.Quantity),
			Currency:    item.UnitPrice.Currency,
		}
		cart.Items = append(cart.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate cart items: %w", err)
	}

	return cart, nil
}

func (r *CartRepository) SaveCart(ctx context.Context, cart *Cart) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	cart.UpdatedAt = time.Now().UTC()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO carts (user_id, updated_at)
		VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET updated_at = EXCLUDED.updated_at
	`, cart.UserID, cart.UpdatedAt); err != nil {
		return fmt.Errorf("failed to upsert cart: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id = $1`, cart.UserID); err != nil {
		return fmt.Errorf("failed to clear cart items: %w", err)
	}

	for _, item := range cart.Items {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO cart_items (
				user_id, product_id, product_name, quantity, unit_price_cents, currency, image_url
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, cart.UserID, item.ProductID, item.ProductName, item.Quantity, item.UnitPrice.AmountCents, item.UnitPrice.Currency, item.ImageURL); err != nil {
			return fmt.Errorf("failed to insert cart item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *CartRepository) DeleteCart(ctx context.Context, userID string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM carts WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	return nil
}

func (r *CartRepository) AddItem(ctx context.Context, userID string, item CartItem) (*Cart, error) {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	found := false
	for i, existing := range cart.Items {
		if existing.ProductID == item.ProductID {
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
