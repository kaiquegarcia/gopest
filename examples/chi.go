package examples

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func customizer(mux *chi.Mux) {
	// Middlewares
	mux.Use(middleware.RequestID)

	// Routes
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := ctx.Value(middleware.RequestIDKey)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"text": "Hello world!",
			"request-id": "%v"
		}`, requestID)))
	})
}
