package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/safar/microservices-demo/gateway/internal/client"
	"github.com/safar/microservices-demo/gateway/internal/errors"
	"github.com/safar/microservices-demo/gateway/internal/middleware"
	"github.com/safar/microservices-demo/gateway/internal/validation"
	cartpb "github.com/safar/microservices-demo/proto/cart/v1"
)

type CartHandler struct {
	cartClient *client.CartClient
}

func NewCartHandler(cartClient *client.CartClient) *CartHandler {
	return &CartHandler{
		cartClient: cartClient,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	resp, err := h.cartClient.GetCart(r.Context(), &cartpb.GetCartRequest{
		UserId: userID,
	})
	if err != nil {
		log.Printf("Failed to get cart for user %s: %v", userID, err)
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	var req cartpb.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode add item request: %v", err)
		errors.WriteError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	req.UserId = userID

	// Validate input
	validationErrors := validation.Validate(
		func() *errors.ValidationError { return validation.ValidateRequired("product_id", req.ProductId) },
		func() *errors.ValidationError { return validation.ValidatePositive("quantity", int64(req.Quantity)) },
		func() *errors.ValidationError {
			if req.UnitPrice == nil {
				return &errors.ValidationError{
					Field:   "unit_price",
					Message: "unit_price is required",
				}
			}
			return nil
		},
	)

	if len(validationErrors) > 0 {
		errors.WriteValidationError(w, validationErrors)
		return
	}

	resp, err := h.cartClient.AddItem(r.Context(), &req)
	if err != nil {
		log.Printf("Failed to add item to cart for user %s: %v", userID, err)
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	productID := chi.URLParam(r, "id")
	if productID == "" {
		errors.WriteError(w, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	var req struct {
		Quantity int32 `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode update item request: %v", err)
		errors.WriteError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate input
	validationErrors := validation.Validate(
		func() *errors.ValidationError { return validation.ValidatePositive("quantity", int64(req.Quantity)) },
	)

	if len(validationErrors) > 0 {
		errors.WriteValidationError(w, validationErrors)
		return
	}

	resp, err := h.cartClient.UpdateItem(r.Context(), &cartpb.UpdateItemRequest{
		UserId:    userID,
		ProductId: productID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		log.Printf("Failed to update cart item for user %s: %v", userID, err)
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	productID := chi.URLParam(r, "id")
	if productID == "" {
		errors.WriteError(w, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	resp, err := h.cartClient.RemoveItem(r.Context(), &cartpb.RemoveItemRequest{
		UserId:    userID,
		ProductId: productID,
	})
	if err != nil {
		log.Printf("Failed to remove cart item for user %s: %v", userID, err)
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	if err := h.cartClient.ClearCart(r.Context(), &cartpb.ClearCartRequest{
		UserId: userID,
	}); err != nil {
		log.Printf("Failed to clear cart for user %s: %v", userID, err)
		errors.WriteGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
