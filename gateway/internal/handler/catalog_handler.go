package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	commonpb "github.com/safar/microservices-demo/proto/common/v1"
	catalogpb "github.com/safar/microservices-demo/proto/catalog/v1"
	"github.com/safar/microservices-demo/gateway/internal/client"
)

type CatalogHandler struct {
	catalogClient *client.CatalogClient
}

func NewCatalogHandler(catalogClient *client.CatalogClient) *CatalogHandler {
	return &CatalogHandler{
		catalogClient: catalogClient,
	}
}

func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	categoryID := r.URL.Query().Get("category_id")
	activeOnly := r.URL.Query().Get("active_only") == "true"

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	resp, err := h.catalogClient.ListProducts(r.Context(), &catalogpb.ListProductsRequest{
		Pagination: &commonpb.Pagination{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
		CategoryId: categoryID,
		ActiveOnly: activeOnly,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := h.catalogClient.GetProduct(r.Context(), &catalogpb.GetProductRequest{
		Identifier: &catalogpb.GetProductRequest_Id{Id: id},
	})
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	categoryID := r.URL.Query().Get("category_id")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	resp, err := h.catalogClient.SearchProducts(r.Context(), &catalogpb.SearchProductsRequest{
		Query: query,
		Pagination: &commonpb.Pagination{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
		CategoryId: categoryID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := h.catalogClient.ListCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req catalogpb.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.catalogClient.CreateProduct(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req catalogpb.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Id = id

	resp, err := h.catalogClient.UpdateProduct(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.catalogClient.DeleteProduct(r.Context(), &catalogpb.DeleteProductRequest{
		Id: id,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
