package client

import (
	"context"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc/metadata"
)

const requestIDHeader = "x-request-id"

func withRequestID(ctx context.Context) context.Context {
	if requestID := chimiddleware.GetReqID(ctx); requestID != "" {
		return metadata.AppendToOutgoingContext(ctx, requestIDHeader, requestID)
	}
	return ctx
}
