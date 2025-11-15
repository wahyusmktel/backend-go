package repository

import (
	"context"
	"database/sql"

	"github.com/nama-github-kamu/bookwise-backend/internal/models"
)

// CategoryRepository mendefinisikan interface untuk CRUD kategori
type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error) // Untuk cek duplikat
}

type mysqlCategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository membuat instance repo kategori
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &mysqlCategoryRepository{db: db}
}

func (r *mysqlCategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := "INSERT INTO categories (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, category.ID, category.Name, category.CreatedAt, category.UpdatedAt)
	return err
}

func (r *mysqlCategoryRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	query := "SELECT id, name, created_at, updated_at FROM categories ORDER BY name ASC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *mysqlCategoryRepository) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	query := "SELECT id, name, created_at, updated_at FROM categories WHERE id = ?"
	var category models.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err // Bisa jadi sql.ErrNoRows
	}
	return &category, nil
}

func (r *mysqlCategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := "UPDATE categories SET name = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, category.Name, category.UpdatedAt, category.ID)
	return err
}

func (r *mysqlCategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	query := "DELETE FROM categories WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *mysqlCategoryRepository) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	query := "SELECT id, name, created_at, updated_at FROM categories WHERE name = ?"
	var category models.Category
	err := r.db.QueryRowContext(ctx, query, name).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err // Pasti sql.ErrNoRows jika tidak ada
	}
	return &category, nil
}