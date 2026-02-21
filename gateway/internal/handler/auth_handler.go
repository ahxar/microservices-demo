package handler

import (
	"encoding/json"
	"log"
	"net/http"

	userpb "github.com/safar/microservices-demo/proto/user/v1"
	"github.com/safar/microservices-demo/gateway/internal/client"
	"github.com/safar/microservices-demo/gateway/internal/errors"
	"github.com/safar/microservices-demo/gateway/internal/validation"
)

type AuthHandler struct {
	userClient *client.UserClient
}

func NewAuthHandler(userClient *client.UserClient) *AuthHandler {
	return &AuthHandler{
		userClient: userClient,
	}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode register request: %v", err)
		errors.WriteError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate input
	validationErrors := validation.Validate(
		func() *errors.ValidationError { return validation.ValidateEmail(req.Email) },
		func() *errors.ValidationError { return validation.ValidatePassword(req.Password) },
		func() *errors.ValidationError { return validation.ValidateRequired("first_name", req.FirstName) },
		func() *errors.ValidationError { return validation.ValidateRequired("last_name", req.LastName) },
	)

	if len(validationErrors) > 0 {
		errors.WriteValidationError(w, validationErrors)
		return
	}

	// Call user service
	resp, err := h.userClient.Register(r.Context(), &userpb.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		errors.WriteGRPCError(w, err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode login request: %v", err)
		errors.WriteError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate input
	validationErrors := validation.Validate(
		func() *errors.ValidationError { return validation.ValidateEmail(req.Email) },
		func() *errors.ValidationError { return validation.ValidateRequired("password", req.Password) },
	)

	if len(validationErrors) > 0 {
		errors.WriteValidationError(w, validationErrors)
		return
	}

	// Call user service
	resp, err := h.userClient.Login(r.Context(), &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Failed to login user: %v", err)
		errors.WriteGRPCError(w, err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode refresh token request: %v", err)
		errors.WriteError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate input
	validationErrors := validation.Validate(
		func() *errors.ValidationError { return validation.ValidateRequired("refresh_token", req.RefreshToken) },
	)

	if len(validationErrors) > 0 {
		errors.WriteValidationError(w, validationErrors)
		return
	}

	// Call user service
	resp, err := h.userClient.RefreshToken(r.Context(), &userpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		errors.WriteGRPCError(w, err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
