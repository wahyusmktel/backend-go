package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10" // validator sudah kita install
	"github.com/nama-github-kamu/bookwise-backend/internal/models"
	"github.com/nama-github-kamu/bookwise-backend/internal/services"
)

// (validator 'validate' sudah ada di auth_handler.go, tapi kita buat lokal saja)
var catValidate = validator.New()

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(s services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	if err := catValidate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	category, err := h.service.CreateCategory(r.Context(), req)
	if err != nil {
		if err.Error() == "nama kategori sudah ada" {
			RespondWithError(w, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Ambil 'id' dari URL

	category, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		if err.Error() == "kategori tidak ditemukan" {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	if err := catValidate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Data input tidak valid")
		return
	}

	category, err := h.service.UpdateCategory(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, errors.New("kategori tidak ditemukan")) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "nama kategori sudah ada" {
			RespondWithError(w, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("kategori tidak ditemukan")) {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Kirim response 204 No Content (standar untuk DELETE sukses)
	w.WriteHeader(http.StatusNoContent)
}