package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nama-github-kamu/bookwise-backend/internal/handlers"
)

// RegisterAuthRoutes mendaftarkan rute-rute autentikasi
func RegisterAuthRoutes(r *chi.Mux, authHandler *handlers.AuthHandler) {
	
    // Grup rute di bawah /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})
}