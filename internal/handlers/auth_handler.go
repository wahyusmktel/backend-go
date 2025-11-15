package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nama-github-kamu/bookwise-backend/internal/models"
	"github.com/nama-github-kamu/bookwise-backend/internal/services"
)

// validate adalah instance global dari validator
var validate = validator.New()

// AuthHandler adalah struct yang menampung dependency untuk auth handler
type AuthHandler struct {
	authService services.AuthService // Dependency ke Service
}

// NewAuthHandler membuat instance AuthHandler baru
func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

// === HANDLER UNTUK REGISTER ===
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	// 1. Decode request JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	// 2. Validasi struct
	if err := validate.Struct(req); err != nil {
		// Mengambil error validasi
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			log.Printf("Error validasi: %v", validationErrors)
			respondWithError(w, http.StatusBadRequest, "Data input tidak valid")
			return
		}
		respondWithError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	// 3. Panggil Service
	userResponse, err := h.authService.RegisterUser(r.Context(), req)
	if err != nil {
		// Cek error spesifik dari service
		if err.Error() == "email sudah terdaftar" {
			respondWithError(w, http.StatusConflict, "Email sudah terdaftar")
			return
		}
		// Error server lainnya
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Kirim response sukses
	respondWithJSON(w, http.StatusCreated, userResponse)
}

// === HANDLER UNTUK LOGIN ===
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	// 1. Decode request JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	// 2. Validasi struct
	if err := validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	// 3. Panggil Service
	loginResponse, err := h.authService.LoginUser(r.Context(), req)
	if err != nil {
		// Cek error spesifik dari service
		if err.Error() == "email atau password salah" {
			// SANGAT PENTING: Gunakan 401 Unauthorized untuk login gagal
			respondWithError(w, http.StatusUnauthorized, "Email atau password salah")
			return
		}
		// Error server lainnya
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Kirim response sukses (berisi user dan token)
	respondWithJSON(w, http.StatusOK, loginResponse)
}

// === FUNGSI HELPER UNTUK RESPONSE ===

// respondWithError mengirim response error JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON mengirim response JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}