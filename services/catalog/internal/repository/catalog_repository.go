package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Category struct {
	ID          string
	Name        string
	Slug        string
	Description string
	ParentID    sql.NullString
	CreatedAt   time.Time
}

type Product struct {
	ID            string
	Name          string
	Slug          string
	Description   string
	PriceCents    int64
	Currency      string
	CategoryID    sql.NullString
	ImageURLs     []string
	StockQuantity int32
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type InventoryReservation struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int32
	ReservedAt time.Time
	ExpiresAt  time.Time
}

type CatalogRepository struct {
	db *sql.DB
}

func NewCatalogRepository(databaseURL string) (*CatalogRepository, error) {
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

	return &CatalogRepository{db: db}, nil
}

func (r *CatalogRepository) Close() error {
	return r.db.Close()
}

// Category operations
func (r *CatalogRepository) ListCategories() ([]*Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, created_at
		FROM categories
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		cat := &Category{}
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.ParentID, &cat.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

// Product operations
func (r *CatalogRepository) CreateProduct(name, slug, description string, priceCents int64, currency, categoryID string, imageURLs []string, stockQuantity int32) (*Product, error) {
	query := `
		INSERT INTO products (name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity, is_active, created_at, updated_at
	`

	var categoryIDNull sql.NullString
	if categoryID != "" {
		categoryIDNull = sql.NullString{String: categoryID, Valid: true}
	}

	product := &Product{}
	err := r.db.QueryRow(query, name, slug, description, priceCents, currency, categoryIDNull, pq.Array(imageURLs), stockQuantity).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description, &product.PriceCents,
		&product.Currency, &product.CategoryID, pq.Array(&product.ImageURLs), &product.StockQuantity,
		&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (r *CatalogRepository) GetProductByID(id string) (*Product, error) {
	query := `
		SELECT id, name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &Product{}
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description, &product.PriceCents,
		&product.Currency, &product.CategoryID, pq.Array(&product.ImageURLs), &product.StockQuantity,
		&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (r *CatalogRepository) GetProductBySlug(slug string) (*Product, error) {
	query := `
		SELECT id, name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity, is_active, created_at, updated_at
		FROM products
		WHERE slug = $1
	`

	product := &Product{}
	err := r.db.QueryRow(query, slug).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description, &product.PriceCents,
		&product.Currency, &product.CategoryID, pq.Array(&product.ImageURLs), &product.StockQuantity,
		&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (r *CatalogRepository) ListProducts(limit, offset int, categoryID string, activeOnly bool) ([]*Product, int, error) {
	countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`
	query := `
		SELECT id, name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity, is_active, created_at, updated_at
		FROM products
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if categoryID != "" {
		countQuery += fmt.Sprintf(" AND category_id = $%d", argPos)
		query += fmt.Sprintf(" AND category_id = $%d", argPos)
		args = append(args, categoryID)
		argPos++
	}

	if activeOnly {
		countQuery += " AND is_active = true"
		query += " AND is_active = true"
	}

	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		product := &Product{}
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.PriceCents,
			&product.Currency, &product.CategoryID, pq.Array(&product.ImageURLs), &product.StockQuantity,
			&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, totalCount, nil
}

func (r *CatalogRepository) SearchProducts(searchQuery string, limit, offset int, categoryID string) ([]*Product, int, error) {
	countQuery := `
		SELECT COUNT(*) FROM products
		WHERE (name ILIKE $1 OR description ILIKE $1)
	`
	query := `
		SELECT id, name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity, is_active, created_at, updated_at
		FROM products
		WHERE (name ILIKE $1 OR description ILIKE $1)
	`

	searchPattern := "%" + searchQuery + "%"
	args := []interface{}{searchPattern}
	argPos := 2

	if categoryID != "" {
		countQuery += fmt.Sprintf(" AND category_id = $%d", argPos)
		query += fmt.Sprintf(" AND category_id = $%d", argPos)
		args = append(args, categoryID)
		argPos++
	}

	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		product := &Product{}
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Description, &product.PriceCents,
			&product.Currency, &product.CategoryID, pq.Array(&product.ImageURLs), &product.StockQuantity,
			&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, totalCount, nil
}

func (r *CatalogRepository) UpdateProduct(id, name, slug, description string, priceCents int64, currency, categoryID string, imageURLs []string, stockQuantity int32, isActive bool) (*Product, error) {
	query := `
		UPDATE products
		SET name = $2, slug = $3, description = $4, price_cents = $5, currency = $6, category_id = $7, image_urls = $8, stock_quantity = $9, is_active = $10, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, slug, description, price_cents, currency, category_id, image_urls, stock_quantity, is_active, created_at, updated_at
	`

	var categoryIDNull sql.NullString
	if categoryID != "" {
		categoryIDNull = sql.NullString{String: categoryID, Valid: true}
	}

	product := &Product{}
	err := r.db.QueryRow(query, id, name, slug, description, priceCents, currency, categoryIDNull, pq.Array(imageURLs), stockQuantity, isActive).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description, &product.PriceCents,
		&product.Currency, &product.CategoryID, pq.Array(&product.ImageURLs), &product.StockQuantity,
		&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (r *CatalogRepository) DeleteProduct(id string) error {
	query := `UPDATE products SET is_active = false WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// Inventory operations
func (r *CatalogRepository) CheckInventory(productID string, quantity int32) (bool, error) {
	query := `SELECT stock_quantity FROM products WHERE id = $1`

	var stockQuantity int32
	if err := r.db.QueryRow(query, productID).Scan(&stockQuantity); err != nil {
		return false, fmt.Errorf("failed to check inventory: %w", err)
	}

	return stockQuantity >= quantity, nil
}

func (r *CatalogRepository) ReserveInventory(orderID, productID string, quantity int32, expirationMinutes int32) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check current stock
	var stockQuantity int32
	if err := tx.QueryRow(`SELECT stock_quantity FROM products WHERE id = $1 FOR UPDATE`, productID).Scan(&stockQuantity); err != nil {
		return "", fmt.Errorf("failed to check stock: %w", err)
	}

	if stockQuantity < quantity {
		return "", fmt.Errorf("insufficient stock")
	}

	// Update stock
	_, err = tx.Exec(`UPDATE products SET stock_quantity = stock_quantity - $1 WHERE id = $2`, quantity, productID)
	if err != nil {
		return "", fmt.Errorf("failed to update stock: %w", err)
	}

	// Create reservation
	expiresAt := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	var reservationID string
	err = tx.QueryRow(`
		INSERT INTO inventory_reservations (order_id, product_id, quantity, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, orderID, productID, quantity, expiresAt).Scan(&reservationID)
	if err != nil {
		return "", fmt.Errorf("failed to create reservation: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return reservationID, nil
}
