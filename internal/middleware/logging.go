package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Middleware func(next http.Handler) http.Handler

func LoggingMW(logger *zap.SugaredLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := uuid.NewString()

			logger.Infow(
				"Got request: ",
				zap.String("ID", reqID),
				zap.String("method", r.Method),
				zap.String("host", r.Host),
				zap.String("proto", r.Proto),
				zap.String("path", r.URL.String()),
				zap.String("content-type", r.Header.Get("Content-Type")),
				zap.Int64("content_length", r.ContentLength),
				zap.String("address", r.RemoteAddr),
			)

			defer func(t time.Time) {
				logger.Infow(
					"Responsed:",
					zap.String("ID", reqID),
					zap.Int64("processing_time_us", time.Since(t).Microseconds()),
					zap.String("content-type", w.Header().Get("Content-Type")),
				)
			}(time.Now())

			next.ServeHTTP(w, r)
		})
	}
}
