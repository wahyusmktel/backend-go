package repository

import (
	"context"
	"database/sql"

	"github.com/nama-github-kamu/bookwise-backend/internal/models"
)

// UserRepository mendefinisikan interface untuk operasi database user.
// Menggunakan interface adalah best practice untuk testing (mocking).
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

// mysqlUserRepository adalah implementasi UserRepository untuk MySQL
type mysqlUserRepository struct {
	db *sql.DB // Dependency ke database
}

// NewUserRepository membuat instance baru dari mysqlUserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

// CreateUser menyimpan user baru ke database
func (r *mysqlUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Query ini harus sesuai dengan skema tabel 'users' kita
	query := "INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"

	// db.ExecContext digunakan untuk query yang tidak mengembalikan baris (INSERT, UPDATE, DELETE)
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password, // Ini adalah password yang SUDAH di-hash
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetUserByEmail mencari user berdasarkan email
func (r *mysqlUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?"

	user := &models.User{}

	// db.QueryRowContext digunakan untuk query yang HANYA mengembalikan 1 baris
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// 'err' di sini bisa jadi 'sql.ErrNoRows' (user tidak ditemukan)
		// atau error koneksi lainnya. Service kita akan menanganinya.
		return nil, err
	}

	return user, nil
}