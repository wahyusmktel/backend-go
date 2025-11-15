package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	// Import semua package kita
	"github.com/nama-github-kamu/bookwise-backend/internal/database"
	"github.com/nama-github-kamu/bookwise-backend/internal/handlers"
	"github.com/nama-github-kamu/bookwise-backend/internal/repository"
	"github.com/nama-github-kamu/bookwise-backend/internal/routes"
	"github.com/nama-github-kamu/bookwise-backend/internal/services"
)

func main() {
	// 1. Muat .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak bisa memuat file .env")
	}

	// 2. Inisialisasi Database
	db := database.InitDB()
	defer db.Close()

	// 3. Setup Router (Chi)
	r := chi.NewRouter()

	// 4. Setup Middleware (CORS)
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:5173"} // Default
	}

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler)

	// 5. ✨ INISIALISASI SEMUA LAYER (DEPENDENCY INJECTION) ✨
	// Buat Repository (bergantung pada DB)
	userRepo := repository.NewUserRepository(db)

	// Buat Service (bergantung pada Repository)
	authService := services.NewAuthService(userRepo)

	// Buat Handler (bergantung pada Service)
	authHandler := handlers.NewAuthHandler(authService)

	// 6. ✨ DAFTARKAN RUTE ✨
	// (bergantung pada Handler)
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	routes.RegisterAuthRoutes(r, authHandler)
	routes.RegisterCategoryRoutes(r, categoryHandler)

	// 7. Rute Health Check (opsional tapi bagus)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "BookWise API is running!"})
	})

	// 8. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("FATAL: Gagal menjalankan server: %v", err)
	}
}

// Helper kecil untuk rute health check di main
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}