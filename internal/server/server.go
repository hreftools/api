package server

import (
	"github.com/zapi-sh/api/internal/handlers"
	"github.com/zapi-sh/api/internal/middlewares"
	"net/http"
	"os"
	"time"
)

func New() *http.Server {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", handlers.Status)
	mux.HandleFunc("GET /resources", handlers.ResourcesGet)

	// version api
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	// apply middlewares
	middlewaresStack := middlewares.MiddlewareStac(
		middlewares.Logging,
		middlewares.CommonHeaders,
	)

	return &http.Server{
		Addr:              ":" + port,
		Handler:           middlewaresStack(v1),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

}
