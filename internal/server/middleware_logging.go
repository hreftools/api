package server

import (
	"log/slog"
	"net/http"
	"time"
)

type wrapperWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapperWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrapperWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		// Map status code to log level so platforms (Railway, etc.) classify
		// 5xx as errors and 4xx as warnings instead of bucketing everything
		// the same.
		level := slog.LevelInfo
		switch {
		case wrapped.statusCode >= 500:
			level = slog.LevelError
		case wrapped.statusCode >= 400:
			level = slog.LevelWarn
		}

		slog.Log(r.Context(), level, "http request",
			"status", wrapped.statusCode,
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
