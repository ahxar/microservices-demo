package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Order struct {
	ID               string
	UserID           string
	Status           string
	SubtotalCents    int64
	ShippingCents    int64
	TaxCents         int64
	TotalCents       int64
	Currency         string
	ShippingStreet   string
	ShippingCity     string
	ShippingState    string
	ShippingZip      string
	ShippingCountry  string
	PaymentMethodID  sql.NullString
	TransactionID    sql.NullString
	TrackingNumber   sql.NullString
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type OrderItem struct {
	ID             string
	OrderID        string
	ProductID      string
	ProductName    string
	Quantity       int32
	UnitPriceCents int64
	TotalPriceCents int64
}

type OrderStatusHistory struct {
	ID        string
	OrderID   string
	Status    string
	Notes     string
	CreatedAt time.Time
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(databaseURL string) (*OrderRepository, error) {
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

	return &OrderRepository{db: db}, nil
}

func (r *OrderRepository) Close() error {
	return r.db.Close()
}

func (r *OrderRepository) CreateOrder(order *Order, items []OrderItem) (*Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create order
	query := `
		INSERT INTO orders (user_id, status, subtotal_cents, shipping_cents, tax_cents, total_cents, currency,
			shipping_street, shipping_city, shipping_state, shipping_zip, shipping_country,
			payment_method_id, transaction_id, tracking_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, user_id, status, subtotal_cents, shipping_cents, tax_cents, total_cents, currency,
			shipping_street, shipping_city, shipping_state, shipping_zip, shipping_country,
			payment_method_id, transaction_id, tracking_number, created_at, updated_at
	`

	err = tx.QueryRow(query,
		order.UserID, order.Status, order.SubtotalCents, order.ShippingCents, order.TaxCents,
		order.TotalCents, order.Currency, order.ShippingStreet, order.ShippingCity,
		order.ShippingState, order.ShippingZip, order.ShippingCountry,
		order.PaymentMethodID, order.TransactionID, order.TrackingNumber,
	).Scan(
		&order.ID, &order.UserID, &order.Status, &order.SubtotalCents, &order.ShippingCents,
		&order.TaxCents, &order.TotalCents, &order.Currency, &order.ShippingStreet,
		&order.ShippingCity, &order.ShippingState, &order.ShippingZip, &order.ShippingCountry,
		&order.PaymentMethodID, &order.TransactionID, &order.TrackingNumber,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, item := range items {
		itemQuery := `
			INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price_cents, total_price_cents)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.Exec(itemQuery, order.ID, item.ProductID, item.ProductName, item.Quantity, item.UnitPriceCents, item.TotalPriceCents)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Create initial status history
	historyQuery := `
		INSERT INTO order_status_history (order_id, status, notes)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(historyQuery, order.ID, order.Status, "Order created")
	if err != nil {
		return nil, fmt.Errorf("failed to create status history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

func (r *OrderRepository) GetOrder(orderID, userID string) (*Order, []OrderItem, []OrderStatusHistory, error) {
	// Get order
	orderQuery := `
		SELECT id, user_id, status, subtotal_cents, shipping_cents, tax_cents, total_cents, currency,
			shipping_street, shipping_city, shipping_state, shipping_zip, shipping_country,
			payment_method_id, transaction_id, tracking_number, created_at, updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`

	order := &Order{}
	err := r.db.QueryRow(orderQuery, orderID, userID).Scan(
		&order.ID, &order.UserID, &order.Status, &order.SubtotalCents, &order.ShippingCents,
		&order.TaxCents, &order.TotalCents, &order.Currency, &order.ShippingStreet,
		&order.ShippingCity, &order.ShippingState, &order.ShippingZip, &order.ShippingCountry,
		&order.PaymentMethodID, &order.TransactionID, &order.TrackingNumber,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get order items
	itemsQuery := `
		SELECT id, order_id, product_id, product_name, quantity, unit_price_cents, total_price_cents
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(itemsQuery, orderID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Quantity, &item.UnitPriceCents, &item.TotalPriceCents); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	// Get status history
	historyQuery := `
		SELECT id, order_id, status, notes, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	historyRows, err := r.db.Query(historyQuery, orderID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get status history: %w", err)
	}
	defer historyRows.Close()

	var history []OrderStatusHistory
	for historyRows.Next() {
		var h OrderStatusHistory
		if err := historyRows.Scan(&h.ID, &h.OrderID, &h.Status, &h.Notes, &h.CreatedAt); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to scan status history: %w", err)
		}
		history = append(history, h)
	}

	return order, items, history, nil
}

func (r *OrderRepository) ListOrders(userID string, limit, offset int, statusFilter string) ([]*Order, int, error) {
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	query := `
		SELECT id, user_id, status, subtotal_cents, shipping_cents, tax_cents, total_cents, currency,
			shipping_street, shipping_city, shipping_state, shipping_zip, shipping_country,
			payment_method_id, transaction_id, tracking_number, created_at, updated_at
		FROM orders
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argPos := 2

	if statusFilter != "" {
		countQuery += fmt.Sprintf(" AND status = $%d", argPos)
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, statusFilter)
		argPos++
	}

	var totalCount int
	if err := r.db.QueryRow(countQuery, args[:len(args)-1]...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		if err := rows.Scan(
			&order.ID, &order.UserID, &order.Status, &order.SubtotalCents, &order.ShippingCents,
			&order.TaxCents, &order.TotalCents, &order.Currency, &order.ShippingStreet,
			&order.ShippingCity, &order.ShippingState, &order.ShippingZip, &order.ShippingCountry,
			&order.PaymentMethodID, &order.TransactionID, &order.TrackingNumber,
			&order.CreatedAt, &order.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, totalCount, nil
}

func (r *OrderRepository) UpdateOrderStatus(orderID, status, notes string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update order status
	_, err = tx.Exec(`UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Add status history
	_, err = tx.Exec(`INSERT INTO order_status_history (order_id, status, notes) VALUES ($1, $2, $3)`, orderID, status, notes)
	if err != nil {
		return fmt.Errorf("failed to add status history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
