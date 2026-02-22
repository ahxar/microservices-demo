package observability

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// StartMetricsServer starts an HTTP server exposing Prometheus metrics and health endpoints.
func StartMetricsServer(logger *slog.Logger, serviceName, metricsPort string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{Status: "ok", Service: serviceName})
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{Status: "ready", Service: serviceName})
	})

	srv := &http.Server{
		Addr:    ":" + metricsPort,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("metrics server failed", slog.String("error", err.Error()), slog.String("addr", srv.Addr))
		}
	}()

	return srv
}

func ShutdownServer(ctx context.Context, logger *slog.Logger, srv *http.Server, name string) {
	if srv == nil {
		return
	}
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", slog.String("server", name), slog.String("error", err.Error()))
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload healthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
