package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blackmamoth/cloudmesh/pkg/config"
	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

func HttpRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		scheme := "http"

		if r.TLS != nil {
			scheme = "https"
		}

		URL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.String())

		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("url", URL),
			zap.String("protocol", r.Proto),
			zap.String("remote_address", r.RemoteAddr),
			zap.Int("status_code", lrw.statusCode),
			zap.Int("size", lrw.size),
			zap.Duration("duration", duration),
		}

		if lrw.statusCode >= http.StatusInternalServerError {
			config.LOGGER.Error("HTTP request error", fields...)
		} else if lrw.statusCode >= http.StatusBadRequest {
			config.LOGGER.Warn("HTTP client error", fields...)
		} else {
			config.LOGGER.Info("HTTP request", fields...)
		}
	})
}
