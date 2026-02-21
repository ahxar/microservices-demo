package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Profile struct {
	ID        string
	UserID    string
	FirstName string
	LastName  string
	Phone     string
	AvatarURL string
}

type Address struct {
	ID        string
	UserID    string
	Label     string
	Street    string
	City      string
	State     string
	ZipCode   string
	Country   string
	IsDefault bool
}

type WishlistItem struct {
	ID        string
	UserID    string
	ProductID string
	AddedAt   time.Time
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(databaseURL string) (*UserRepository, error) {
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

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Close() error {
	return r.db.Close()
}

// User operations
func (r *UserRepository) CreateUser(email, passwordHash, role string) (*User, error) {
	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role, created_at, updated_at
	`

	user := &User{}
	err := r.db.QueryRow(query, email, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *UserRepository) ListUsers(limit, offset int) ([]*User, int, error) {
	countQuery := `SELECT COUNT(*) FROM users`
	var totalCount int
	if err := r.db.QueryRow(countQuery).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Role,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, totalCount, nil
}

// Profile operations
func (r *UserRepository) CreateProfile(userID, firstName, lastName, phone, avatarURL string) (*Profile, error) {
	query := `
		INSERT INTO profiles (user_id, first_name, last_name, phone, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, first_name, last_name, phone, avatar_url
	`

	profile := &Profile{}
	err := r.db.QueryRow(query, userID, firstName, lastName, phone, avatarURL).Scan(
		&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName,
		&profile.Phone, &profile.AvatarURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return profile, nil
}

func (r *UserRepository) GetProfileByUserID(userID string) (*Profile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, avatar_url
		FROM profiles
		WHERE user_id = $1
	`

	profile := &Profile{}
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName,
		&profile.Phone, &profile.AvatarURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return profile, nil
}

func (r *UserRepository) UpdateProfile(userID, firstName, lastName, phone, avatarURL string) (*Profile, error) {
	query := `
		UPDATE profiles
		SET first_name = $2, last_name = $3, phone = $4, avatar_url = $5
		WHERE user_id = $1
		RETURNING id, user_id, first_name, last_name, phone, avatar_url
	`

	profile := &Profile{}
	err := r.db.QueryRow(query, userID, firstName, lastName, phone, avatarURL).Scan(
		&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName,
		&profile.Phone, &profile.AvatarURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return profile, nil
}

// Address operations
func (r *UserRepository) CreateAddress(userID, label, street, city, state, zipCode, country string, isDefault bool) (*Address, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if isDefault {
		_, err = tx.Exec(`UPDATE addresses SET is_default = false WHERE user_id = $1`, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to reset default addresses: %w", err)
		}
	}

	query := `
		INSERT INTO addresses (user_id, label, street, city, state, zip_code, country, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, label, street, city, state, zip_code, country, is_default
	`

	address := &Address{}
	err = tx.QueryRow(query, userID, label, street, city, state, zipCode, country, isDefault).Scan(
		&address.ID, &address.UserID, &address.Label, &address.Street, &address.City,
		&address.State, &address.ZipCode, &address.Country, &address.IsDefault,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return address, nil
}

func (r *UserRepository) ListAddresses(userID string) ([]*Address, error) {
	query := `
		SELECT id, user_id, label, street, city, state, zip_code, country, is_default
		FROM addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}
	defer rows.Close()

	var addresses []*Address
	for rows.Next() {
		address := &Address{}
		if err := rows.Scan(
			&address.ID, &address.UserID, &address.Label, &address.Street, &address.City,
			&address.State, &address.ZipCode, &address.Country, &address.IsDefault,
		); err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Wishlist operations
func (r *UserRepository) AddToWishlist(userID, productID string) error {
	query := `
		INSERT INTO wishlists (user_id, product_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, product_id) DO NOTHING
	`

	_, err := r.db.Exec(query, userID, productID)
	if err != nil {
		return fmt.Errorf("failed to add to wishlist: %w", err)
	}

	return nil
}

func (r *UserRepository) GetWishlist(userID string) ([]*WishlistItem, error) {
	query := `
		SELECT id, user_id, product_id, added_at
		FROM wishlists
		WHERE user_id = $1
		ORDER BY added_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}
	defer rows.Close()

	var items []*WishlistItem
	for rows.Next() {
		item := &WishlistItem{}
		if err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.AddedAt); err != nil {
			return nil, fmt.Errorf("failed to scan wishlist item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
