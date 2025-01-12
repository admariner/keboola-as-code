package middleware

import (
	"context"
	"net/http"

	goaMiddleware "goa.design/goa/v3/middleware"

	"github.com/keboola/keboola-as-code/internal/pkg/idgenerator"
)

const (
	RequestIDHeader  = "X-Request-Id"
	RequestIDCtxKey  = ctxKey("request-id")
	RequestURLCtxKey = ctxKey("request-url")
)

// RequestInfo middleware adds requestID and URL to the context values.
func RequestInfo() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Generate unique request ID
			requestID := idgenerator.RequestID()

			// Update context
			ctx := req.Context()
			ctx = context.WithValue(ctx, goaMiddleware.RequestIDKey, requestID) // nolint:staticcheck // intentionally used the ctx key from external package
			ctx = context.WithValue(ctx, RequestIDCtxKey, requestID)
			ctx = context.WithValue(ctx, RequestURLCtxKey, req.URL)
			req = req.WithContext(ctx)

			// Add request ID to headers
			w.Header().Add(RequestIDHeader, requestID)

			// Handle
			next.ServeHTTP(w, req)
		})
	}
}
