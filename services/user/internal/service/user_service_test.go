package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/safar/microservices-demo/services/user/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	getUserByEmailFn func(email string) (*repository.User, error)
	getProfileByIDFn func(userID string) (*repository.Profile, error)
	listUsersFn      func(limit, offset int) ([]*repository.User, int, error)
}

func (m *mockUserRepository) CreateUser(email, passwordHash, role string) (*repository.User, error) {
	return nil, nil
}

func (m *mockUserRepository) CreateProfile(userID, firstName, lastName, phone, avatarURL string) (*repository.Profile, error) {
	return nil, nil
}

func (m *mockUserRepository) GetUserByEmail(email string) (*repository.User, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(email)
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepository) GetProfileByUserID(userID string) (*repository.Profile, error) {
	if m.getProfileByIDFn != nil {
		return m.getProfileByIDFn(userID)
	}
	return nil, nil
}

func (m *mockUserRepository) GetUserByID(id string) (*repository.User, error) {
	return nil, nil
}

func (m *mockUserRepository) UpdateProfile(userID, firstName, lastName, phone, avatarURL string) (*repository.Profile, error) {
	return nil, nil
}

func (m *mockUserRepository) ListUsers(limit, offset int) ([]*repository.User, int, error) {
	if m.listUsersFn != nil {
		return m.listUsersFn(limit, offset)
	}
	return nil, 0, nil
}

func (m *mockUserRepository) CreateAddress(userID, label, street, city, state, zipCode, country string, isDefault bool) (*repository.Address, error) {
	return nil, nil
}

func (m *mockUserRepository) ListAddresses(userID string) ([]*repository.Address, error) {
	return nil, nil
}

func (m *mockUserRepository) AddToWishlist(userID, productID string) error {
	return nil
}

func (m *mockUserRepository) GetWishlist(userID string) ([]*repository.WishlistItem, error) {
	return nil, nil
}

func TestLoginAndValidateToken(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	mockRepo := &mockUserRepository{
		getUserByEmailFn: func(email string) (*repository.User, error) {
			return &repository.User{
				ID:           "user-1",
				Email:        email,
				PasswordHash: string(hashedPassword),
				Role:         "customer",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
		getProfileByIDFn: func(userID string) (*repository.Profile, error) {
			return &repository.Profile{UserID: userID, FirstName: "Test", LastName: "User"}, nil
		},
	}

	svc := NewUserService(mockRepo, "test-secret", 3600)
	_, _, accessToken, _, err := svc.Login(context.Background(), "user@example.com", "password123")
	if err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}
	if accessToken == "" {
		t.Fatalf("expected access token")
	}

	userID, role, err := svc.ValidateToken(context.Background(), accessToken)
	if err != nil {
		t.Fatalf("expected token validation to succeed, got %v", err)
	}
	if userID != "user-1" || role != "customer" {
		t.Fatalf("unexpected token claims: userID=%s role=%s", userID, role)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailFn: func(email string) (*repository.User, error) {
			return nil, sql.ErrNoRows
		},
	}

	svc := NewUserService(mockRepo, "test-secret", 3600)
	_, _, _, _, err := svc.Login(context.Background(), "missing@example.com", "password123")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestListUsersPaginationOffset(t *testing.T) {
	mockRepo := &mockUserRepository{
		listUsersFn: func(limit, offset int) ([]*repository.User, int, error) {
			if limit != 10 || offset != 20 {
				t.Fatalf("unexpected pagination values: limit=%d offset=%d", limit, offset)
			}
			return []*repository.User{}, 0, nil
		},
	}

	svc := NewUserService(mockRepo, "test-secret", 3600)
	_, _, err := svc.ListUsers(context.Background(), 3, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
