package utils

import (
	"github.com/go-chi/cors"
	"net/http"

	"github.com/go-chi/chi"
)

type Router struct {
	Mux chi.Router
}

func NewRouter() *Router {
	router := Router{
		Mux: chi.NewRouter(),
	}
	return &router
}

type RouterOption func(http.Handler) http.Handler

// UseMiddlewares for registering middleware
func (r *Router) UseMiddlewares(opts ...RouterOption) {
	r.Mux.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)
	for _, opt := range opts {
		r.Mux.Use(opt)
	}
}
