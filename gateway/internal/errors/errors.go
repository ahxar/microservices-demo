package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Code    string            `json:"code,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// WriteError writes a structured error response
func WriteError(w http.ResponseWriter, statusCode int, message string, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Details: details,
		Code:    fmt.Sprintf("%d", statusCode),
	}

	json.NewEncoder(w).Encode(resp)
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, errors []ValidationError) {
	details := make(map[string]string)
	for _, err := range errors {
		details[err.Field] = err.Message
	}

	WriteError(w, http.StatusBadRequest, "Validation failed", details)
}

// GRPCErrorToHTTP converts a gRPC error to appropriate HTTP status code
func GRPCErrorToHTTP(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, "Internal server error"
	}

	switch st.Code() {
	case codes.OK:
		return http.StatusOK, ""
	case codes.InvalidArgument:
		return http.StatusBadRequest, st.Message()
	case codes.NotFound:
		return http.StatusNotFound, st.Message()
	case codes.AlreadyExists:
		return http.StatusConflict, st.Message()
	case codes.PermissionDenied:
		return http.StatusForbidden, st.Message()
	case codes.Unauthenticated:
		return http.StatusUnauthorized, st.Message()
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests, "Rate limit exceeded"
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed, st.Message()
	case codes.Aborted:
		return http.StatusConflict, st.Message()
	case codes.OutOfRange:
		return http.StatusBadRequest, st.Message()
	case codes.Unimplemented:
		return http.StatusNotImplemented, st.Message()
	case codes.Unavailable:
		return http.StatusServiceUnavailable, "Service temporarily unavailable"
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, "Request timeout"
	case codes.Canceled:
		return http.StatusRequestTimeout, "Request canceled"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}

// WriteGRPCError converts and writes a gRPC error as HTTP response
func WriteGRPCError(w http.ResponseWriter, err error) {
	statusCode, message := GRPCErrorToHTTP(err)
	WriteError(w, statusCode, message, nil)
}
