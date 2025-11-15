package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nama-github-kamu/bookwise-backend/internal/models"
	"github.com/nama-github-kamu/bookwise-backend/internal/repository"
)

// CategoryService mendefinisikan interface logika bisnis kategori
type CategoryService interface {
	CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.Category, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
	UpdateCategory(ctx context.Context, id string, req models.UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type categoryServiceImpl struct {
	repo repository.CategoryRepository
}

// NewCategoryService membuat instance service kategori
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{repo: repo}
}

func (s *categoryServiceImpl) CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.Category, error) {
	// Normalisasi & Validasi Logika Bisnis
	name := strings.TrimSpace(req.Name)
	if len(name) < 3 {
		return nil, errors.New("nama kategori terlalu pendek")
	}
	
	// Cek duplikat
	if _, err := s.repo.GetCategoryByName(ctx, name); err == nil {
		return nil, errors.New("nama kategori sudah ada")
	}

	category := &models.Category{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryServiceImpl) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *categoryServiceImpl) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("kategori tidak ditemukan")
		}
		return nil, err
	}
	return category, nil
}

func (s *categoryServiceImpl) UpdateCategory(ctx context.Context, id string, req models.UpdateCategoryRequest) (*models.Category, error) {
	// Cek dulu apakah kategori-nya ada
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("kategori tidak ditemukan")
		}
		return nil, err
	}

	// Normalisasi & Validasi
	name := strings.TrimSpace(req.Name)
	if len(name) < 3 {
		return nil, errors.New("nama kategori terlalu pendek")
	}
	
	// Cek duplikat (pastikan bukan dirinya sendiri)
	existing, err := s.repo.GetCategoryByName(ctx, name)
	if err == nil && existing.ID != id {
		return nil, errors.New("nama kategori sudah ada")
	}

	// Update data
	category.Name = name
	category.UpdatedAt = time.Now()

	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, id string) error {
	// Cek dulu apakah ada
	if _, err := s.repo.GetCategoryByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("kategori tidak ditemukan")
		}
		return err
	}

	// TODO: Cek apakah kategori ini dipakai oleh buku?
	// Jika iya, harusnya kita larang hapus.
	// Untuk sekarang, kita hajar hapus.
	
	return s.repo.DeleteCategory(ctx, id)
}