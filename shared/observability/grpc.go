package observability

import (
	"context"
	"log/slog"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcHealth "google.golang.org/grpc/health"
	grpcHealthV1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func NewGRPCServer(logger *slog.Logger) (*grpc.Server, *grpcprom.ServerMetrics, *grpcHealth.Server) {
	metrics := grpcprom.NewServerMetrics()

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			metrics.UnaryServerInterceptor(),
			UnaryServerLoggingInterceptor(logger),
		),
		grpc.ChainStreamInterceptor(
			metrics.StreamServerInterceptor(),
			StreamServerLoggingInterceptor(logger),
		),
	)

	healthServer := grpcHealth.NewServer()
	grpcHealthV1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpcHealthV1.HealthCheckResponse_SERVING)

	return s, metrics, healthServer
}

func UnaryServerLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)

		logger.Info("grpc request",
			slog.String("transport", "grpc"),
			slog.String("method", info.FullMethod),
			slog.String("peer", peerAddress(ctx)),
			slog.String("code", code.String()),
			slog.Bool("error", code != codes.OK),
			slog.Duration("duration", time.Since(start)),
			slog.String("trace_id", traceIDFromContext(ctx)),
		)

		return resp, err
	}
}

func StreamServerLoggingInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		err := handler(srv, ss)
		code := status.Code(err)
		ctx := ss.Context()

		logger.Info("grpc stream request",
			slog.String("transport", "grpc"),
			slog.String("method", info.FullMethod),
			slog.String("peer", peerAddress(ctx)),
			slog.String("code", code.String()),
			slog.Bool("error", code != codes.OK),
			slog.Duration("duration", time.Since(start)),
			slog.String("trace_id", traceIDFromContext(ctx)),
		)

		return err
	}
}

func peerAddress(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p.Addr == nil {
		return ""
	}
	return p.Addr.String()
}

func traceIDFromContext(ctx context.Context) string {
	sc := trace.SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return ""
	}
	return sc.TraceID().String()
}
