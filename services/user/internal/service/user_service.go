package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/safar/microservices-demo/services/user/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewUserService(repo *repository.UserRepository, jwtSecret string, jwtExpiry int) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// JWT Claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Register creates a new user
func (s *UserService) Register(ctx context.Context, email, password, firstName, lastName string) (*repository.User, *repository.Profile, string, string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.repo.CreateUser(email, string(hashedPassword), "customer")
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to create user: %w", err)
	}

	// Create profile
	profile, err := s.repo.CreateProfile(user.ID, firstName, lastName, "", "")
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to create profile: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, profile, accessToken, refreshToken, nil
}

// Login authenticates a user
func (s *UserService) Login(ctx context.Context, email, password string) (*repository.User, *repository.Profile, string, string, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, "", "", errors.New("invalid credentials")
		}
		return nil, nil, "", "", fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, "", "", errors.New("invalid credentials")
	}

	// Get profile
	profile, err := s.repo.GetProfileByUserID(user.ID)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to get profile: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, profile, accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token
func (s *UserService) ValidateToken(ctx context.Context, tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	return claims.UserID, claims.Role, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*repository.User, *repository.Profile, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	profile, err := s.repo.GetProfileByUserID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return user, profile, nil
}

// UpdateUser updates user profile
func (s *UserService) UpdateUser(ctx context.Context, userID, firstName, lastName, phone, avatarURL string) (*repository.User, *repository.Profile, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if profile exists
	existingProfile, err := s.repo.GetProfileByUserID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get profile: %w", err)
	}

	var profile *repository.Profile
	if existingProfile == nil {
		// Create profile if it doesn't exist
		profile, err = s.repo.CreateProfile(userID, firstName, lastName, phone, avatarURL)
	} else {
		// Update existing profile
		profile, err = s.repo.UpdateProfile(userID, firstName, lastName, phone, avatarURL)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to save profile: %w", err)
	}

	return user, profile, nil
}

// ListUsers lists all users (admin only)
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]*repository.User, int, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.repo.ListUsers(pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// AddAddress adds a new address for a user
func (s *UserService) AddAddress(ctx context.Context, userID, label, street, city, state, zipCode, country string, isDefault bool) (*repository.Address, error) {
	address, err := s.repo.CreateAddress(userID, label, street, city, state, zipCode, country, isDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return address, nil
}

// ListAddresses lists all addresses for a user
func (s *UserService) ListAddresses(ctx context.Context, userID string) ([]*repository.Address, error) {
	addresses, err := s.repo.ListAddresses(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	return addresses, nil
}

// AddToWishlist adds a product to user's wishlist
func (s *UserService) AddToWishlist(ctx context.Context, userID, productID string) error {
	if err := s.repo.AddToWishlist(userID, productID); err != nil {
		return fmt.Errorf("failed to add to wishlist: %w", err)
	}

	return nil
}

// GetWishlist retrieves user's wishlist
func (s *UserService) GetWishlist(ctx context.Context, userID string) ([]*repository.WishlistItem, error) {
	items, err := s.repo.GetWishlist(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	return items, nil
}

// Helper function to generate JWT tokens
func (s *UserService) generateTokens(user *repository.User) (string, string, error) {
	// Access token (short-lived)
	accessClaims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token (long-lived)
	refreshClaims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
