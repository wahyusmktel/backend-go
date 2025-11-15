package models

import "time"

// Category adalah model untuk tabel 'categories' di DB
type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCategoryRequest adalah payload JSON untuk membuat kategori
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

// UpdateCategoryRequest adalah payload JSON untuk update kategori
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}