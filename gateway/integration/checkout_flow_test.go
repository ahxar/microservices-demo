package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

type authResponse struct {
	AccessToken string `json:"access_token"`
}

type money struct {
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

type product struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StockQuantity int32  `json:"stock_quantity"`
	Price         money  `json:"price"`
}

type listProductsResponse struct {
	Products []product `json:"products"`
}

func TestCheckoutFlow(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION") != "1" {
		t.Skip("set RUN_INTEGRATION=1 to run integration tests")
	}

	baseURL := os.Getenv("INTEGRATION_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	email := fmt.Sprintf("integration-%d@example.com", time.Now().UnixNano())

	registerBody := map[string]any{
		"email":      email,
		"password":   "password123",
		"first_name": "Integration",
		"last_name":  "Test",
	}

	registerPayload, _ := json.Marshal(registerBody)
	registerResp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(registerPayload))
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	defer registerResp.Body.Close()
	if registerResp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected register status: %d", registerResp.StatusCode)
	}

	var auth authResponse
	if err := json.NewDecoder(registerResp.Body).Decode(&auth); err != nil {
		t.Fatalf("failed to decode register response: %v", err)
	}
	if auth.AccessToken == "" {
		t.Fatalf("expected access token")
	}

	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/products?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	productsResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to fetch products: %v", err)
	}
	defer productsResp.Body.Close()
	if productsResp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected products status: %d", productsResp.StatusCode)
	}

	var products listProductsResponse
	if err := json.NewDecoder(productsResp.Body).Decode(&products); err != nil {
		t.Fatalf("failed to decode products response: %v", err)
	}
	if len(products.Products) == 0 {
		t.Skip("no products available in catalog seed data")
	}

	selected := products.Products[0]
	addCartBody := map[string]any{
		"product_id":   selected.ID,
		"product_name": selected.Name,
		"quantity":     1,
		"unit_price": map[string]any{
			"amount_cents": selected.Price.AmountCents,
			"currency":     selected.Price.Currency,
		},
	}
	addCartPayload, _ := json.Marshal(addCartBody)
	addReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/cart/items", bytes.NewReader(addCartPayload))
	addReq.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	addReq.Header.Set("Content-Type", "application/json")
	addResp, err := client.Do(addReq)
	if err != nil {
		t.Fatalf("failed to add cart item: %v", err)
	}
	defer addResp.Body.Close()
	if addResp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected add-to-cart status: %d", addResp.StatusCode)
	}

	checkoutBody := map[string]any{
		"shipping_address": map[string]any{
			"street":   "1 Main St",
			"city":     "San Francisco",
			"state":    "CA",
			"zip_code": "94105",
			"country":  "USA",
		},
		// A seed payment method UUID can be passed with INTEGRATION_PAYMENT_METHOD_ID.
		"payment_method_id": os.Getenv("INTEGRATION_PAYMENT_METHOD_ID"),
	}

	if checkoutBody["payment_method_id"] == "" {
		t.Skip("set INTEGRATION_PAYMENT_METHOD_ID to run checkout creation end-to-end")
	}

	checkoutPayload, _ := json.Marshal(checkoutBody)
	checkoutReq, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/orders", bytes.NewReader(checkoutPayload))
	checkoutReq.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	checkoutReq.Header.Set("Content-Type", "application/json")
	checkoutResp, err := client.Do(checkoutReq)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	defer checkoutResp.Body.Close()

	if checkoutResp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected checkout status: %d", checkoutResp.StatusCode)
	}
}
