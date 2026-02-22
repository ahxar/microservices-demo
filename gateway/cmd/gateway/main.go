package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"github.com/safar/microservices-demo/gateway/internal/client"
	"github.com/safar/microservices-demo/gateway/internal/config"
	"github.com/safar/microservices-demo/gateway/internal/handler"
	"github.com/safar/microservices-demo/gateway/internal/middleware"
	"github.com/safar/microservices-demo/gateway/internal/observability"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize tracing (fail-open).
	shutdownTracing := func(context.Context) error { return nil }
	if shutdown, err := observability.InitTracing(
		context.Background(),
		cfg.OTELServiceName,
		cfg.OTELExporterOTLPEndpoint,
		cfg.OTELExporterOTLPInsecure,
	); err != nil {
		log.Printf("warning: tracing disabled: %v", err)
	} else {
		shutdownTracing = shutdown
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(ctx); err != nil {
			log.Printf("warning: tracer shutdown failed: %v", err)
		}
	}()

	metricsServer := startMetricsServer(cfg.MetricsPort)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := metricsServer.Shutdown(ctx); err != nil {
			log.Printf("warning: metrics server shutdown failed: %v", err)
		}
	}()

	// Initialize gRPC clients
	userClient, err := client.NewUserClient(cfg.UserServiceURL)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userClient.Close()

	catalogClient, err := client.NewCatalogClient(cfg.CatalogServiceURL)
	if err != nil {
		log.Fatalf("Failed to connect to catalog service: %v", err)
	}
	defer catalogClient.Close()

	cartClient, err := client.NewCartClient(cfg.CartServiceURL)
	if err != nil {
		log.Fatalf("Failed to connect to cart service: %v", err)
	}
	defer cartClient.Close()

	orderClient, err := client.NewOrderClient(cfg.OrderServiceURL)
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderClient.Close()

	// Initialize rate limiter
	rateLimiter, err := middleware.NewRateLimiter(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialize rate limiter: %v", err)
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userClient)
	userHandler := handler.NewUserHandler(userClient)
	catalogHandler := handler.NewCatalogHandler(catalogClient)
	cartHandler := handler.NewCartHandler(cartClient)
	orderHandler := handler.NewOrderHandler(orderClient)

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(otelchi.Middleware(cfg.OTELServiceName))
	r.Use(middleware.CORS())
	r.Use(middleware.Logger)
	r.Use(middleware.Metrics())
	r.Use(rateLimiter.Middleware(100)) // 100 requests per minute

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes - Authentication
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.RefreshToken)
		})

		// Public routes - Catalog
		r.Get("/products", catalogHandler.ListProducts)
		r.Get("/products/{id}", catalogHandler.GetProduct)
		r.Get("/products/search", catalogHandler.SearchProducts)
		r.Get("/categories", catalogHandler.ListCategories)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))

			// User routes
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateMe)
			r.Get("/addresses", userHandler.ListAddresses)
			r.Post("/addresses", userHandler.AddAddress)
			r.Get("/wishlist", userHandler.GetWishlist)
			r.Post("/wishlist", userHandler.AddToWishlist)

			// Cart routes
			r.Get("/cart", cartHandler.GetCart)
			r.Post("/cart/items", cartHandler.AddItem)
			r.Put("/cart/items/{id}", cartHandler.UpdateItem)
			r.Delete("/cart/items/{id}", cartHandler.RemoveItem)
			r.Delete("/cart", cartHandler.ClearCart)

			// Order routes
			r.Get("/orders", orderHandler.ListOrders)
			r.Post("/orders", orderHandler.CreateOrder)
			r.Get("/orders/{id}", orderHandler.GetOrder)
			r.Delete("/orders/{id}", orderHandler.CancelOrder)
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))
			r.Use(middleware.AdminOnly)

			// Product management
			r.Post("/admin/products", catalogHandler.CreateProduct)
			r.Put("/admin/products/{id}", catalogHandler.UpdateProduct)
			r.Delete("/admin/products/{id}", catalogHandler.DeleteProduct)
			r.Get("/admin/users", userHandler.ListUsers)
		})
	})

	// Start server
	log.Printf("API Gateway starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startMetricsServer(metricsPort string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	srv := &http.Server{
		Addr:    ":" + metricsPort,
		Handler: mux,
	}

	go func() {
		log.Printf("Gateway metrics server starting on port %s", metricsPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("metrics server failed: %v", err)
		}
	}()

	return srv
}
