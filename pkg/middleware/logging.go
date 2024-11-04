package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// HTTP Logger middleware
func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom response writer to capture status code
			wrw := newWrappedResponseWriter(w)

			// Process request
			next.ServeHTTP(wrw, r)

			// Log the request details
			logger.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int("status", wrw.status),
				zap.Duration("latency", time.Since(start)),
				zap.Int("bytes", wrw.bytesWritten),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}

// gRPC Logger interceptor
func GRPCLogger(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Process request
		resp, err := handler(ctx, req)

		// Log the request details
		logger.Info("gRPC Request",
			zap.String("method", info.FullMethod),
			zap.Duration("latency", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}

// Custom response writer to capture status code and bytes written
type wrappedResponseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func newWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrappedResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}
