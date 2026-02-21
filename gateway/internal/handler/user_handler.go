package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/safar/microservices-demo/gateway/internal/client"
	"github.com/safar/microservices-demo/gateway/internal/middleware"
	commonpb "github.com/safar/microservices-demo/proto/common/v1"
	userpb "github.com/safar/microservices-demo/proto/user/v1"
)

type UserHandler struct {
	userClient *client.UserClient
}

func NewUserHandler(userClient *client.UserClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	resp, err := h.userClient.GetUser(r.Context(), &userpb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var profile userpb.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.userClient.UpdateUser(r.Context(), &userpb.UpdateUserRequest{
		Id:      userID,
		Profile: &profile,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	resp, err := h.userClient.ListAddresses(r.Context(), &userpb.ListAddressesRequest{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type AddAddressRequest struct {
	Label     string            `json:"label"`
	Address   *commonpb.Address `json:"address"`
	IsDefault bool              `json:"is_default"`
}

func (h *UserHandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var req AddAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.userClient.AddAddress(r.Context(), &userpb.AddAddressRequest{
		UserId:    userID,
		Label:     req.Label,
		Address:   req.Address,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	resp, err := h.userClient.GetWishlist(r.Context(), &userpb.GetWishlistRequest{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type AddToWishlistRequest struct {
	ProductID string `json:"product_id"`
}

func (h *UserHandler) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var req AddToWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userClient.AddToWishlist(r.Context(), &userpb.AddToWishlistRequest{
		UserId:    userID,
		ProductId: req.ProductID,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "added to wishlist"}`))
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	resp, err := h.userClient.ListUsers(r.Context(), &userpb.ListUsersRequest{
		Pagination: &commonpb.Pagination{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
