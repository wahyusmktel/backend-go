package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nama-github-kamu/bookwise-backend/internal/handlers"
	"github.com/nama-github-kamu/bookwise-backend/internal/middleware"
)

// RegisterCategoryRoutes mendaftarkan rute CRUD untuk Kategori
func RegisterCategoryRoutes(r *chi.Mux, h *handlers.CategoryHandler) {
	
	// Grup rute di bawah /api/v1
	// Kita akan memproteksi SEMUA rute di dalam grup ini
	r.Route("/api/v1/categories", func(r chi.Router) {
		
		// Gunakan AuthMiddleware kita
		r.Use(middleware.AuthMiddleware)

		r.Post("/", h.CreateCategory)     // POST /api/v1/categories
		r.Get("/", h.GetCategories)       // GET /api/v1/categories
		r.Get("/{id}", h.GetCategoryByID) // GET /api/v1/categories/some-uuid
		r.Put("/{id}", h.UpdateCategory)  // PUT /api/v1/categories/some-uuid
		r.Delete("/{id}", h.DeleteCategory) // DELETE /api/v1/categories/some-uuid
	})
}