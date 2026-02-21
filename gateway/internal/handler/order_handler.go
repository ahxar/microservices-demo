package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/safar/microservices-demo/gateway/internal/client"
	"github.com/safar/microservices-demo/gateway/internal/errors"
	"github.com/safar/microservices-demo/gateway/internal/middleware"
	commonpb "github.com/safar/microservices-demo/proto/common/v1"
	orderpb "github.com/safar/microservices-demo/proto/order/v1"
)

type OrderHandler struct {
	orderClient *client.OrderClient
}

func NewOrderHandler(orderClient *client.OrderClient) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
	}
}

type CreateOrderRequest struct {
	ShippingAddress *commonpb.Address `json:"shipping_address"`
	PaymentMethodID string            `json:"payment_method_id"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.ShippingAddress == nil || req.PaymentMethodID == "" {
		errors.WriteError(w, http.StatusBadRequest, "shipping_address and payment_method_id are required", nil)
		return
	}

	resp, err := h.orderClient.CreateOrder(r.Context(), &orderpb.CreateOrderRequest{
		UserId:          userID,
		ShippingAddress: req.ShippingAddress,
		PaymentMethodId: req.PaymentMethodID,
	})
	if err != nil {
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		errors.WriteError(w, http.StatusBadRequest, "Order ID is required", nil)
		return
	}

	resp, err := h.orderClient.GetOrder(r.Context(), &orderpb.GetOrderRequest{
		Id:     orderID,
		UserId: userID,
	})
	if err != nil {
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	statusFilter := parseOrderStatus(r.URL.Query().Get("status"))

	resp, err := h.orderClient.ListOrders(r.Context(), &orderpb.ListOrdersRequest{
		UserId: userID,
		Pagination: &commonpb.Pagination{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
		StatusFilter: statusFilter,
	})
	if err != nil {
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		errors.WriteError(w, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		errors.WriteError(w, http.StatusBadRequest, "Order ID is required", nil)
		return
	}

	resp, err := h.orderClient.CancelOrder(r.Context(), &orderpb.CancelOrderRequest{
		Id:     orderID,
		UserId: userID,
		Reason: r.URL.Query().Get("reason"),
	})
	if err != nil {
		errors.WriteGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func parseOrderStatus(status string) orderpb.OrderStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending":
		return orderpb.OrderStatus_ORDER_STATUS_PENDING
	case "confirmed":
		return orderpb.OrderStatus_ORDER_STATUS_CONFIRMED
	case "processing":
		return orderpb.OrderStatus_ORDER_STATUS_PROCESSING
	case "shipped":
		return orderpb.OrderStatus_ORDER_STATUS_SHIPPED
	case "delivered":
		return orderpb.OrderStatus_ORDER_STATUS_DELIVERED
	case "cancelled":
		return orderpb.OrderStatus_ORDER_STATUS_CANCELLED
	case "refunded":
		return orderpb.OrderStatus_ORDER_STATUS_REFUNDED
	default:
		return orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
