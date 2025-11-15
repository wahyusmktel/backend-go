package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/nama-github-kamu/bookwise-backend/internal/models"
	"github.com/nama-github-kamu/bookwise-backend/internal/repository"
)

// AuthService mendefinisikan interface untuk logika autentikasi
type AuthService interface {
	RegisterUser(ctx context.Context, req models.RegisterRequest) (*models.UserResponse, error)
	LoginUser(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

// authService adalah implementasi dari AuthService
type authService struct {
	userRepo  repository.UserRepository // Dependency ke Repository
	jwtSecret string                    // Kunci rahasIA JWT
}

// NewAuthService membuat instance authService baru
func NewAuthService(repo repository.UserRepository) AuthService {
	// Ambil JWT Secret dari .env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("FATAL: JWT_SECRET tidak diset di .env")
	}

	return &authService{
		userRepo:  repo,
		jwtSecret: secret,
	}
}

// === Implementasi Logika REGISTER ===
func (s *authService) RegisterUser(ctx context.Context, req models.RegisterRequest) (*models.UserResponse, error) {
	// 1. Cek apakah email sudah terdaftar
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		// Ini adalah error database sungguhan, bukan "user tidak ditemukan"
		log.Printf("Error cek email: %v", err)
		return nil, errors.New("terjadi kesalahan pada server")
	}
	if err == nil {
		// User ditemukan (err == nil), artinya email sudah dipakai
		return nil, errors.New("email sudah terdaftar")
	}

	// 2. Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("terjadi kesalahan pada server")
	}

	// 3. Buat struct User baru untuk disimpan
	newUser := &models.User{
		ID:        uuid.NewString(), // Generate UUID baru
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Simpan user ke database
	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		log.Printf("Error menyimpan user: %v", err)
		return nil, errors.New("gagal mendaftarkan user")
	}

	// 5. Kembalikan UserResponse (tanpa password)
	return &models.UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil
}

// === Implementasi Logika LOGIN ===
func (s *authService) LoginUser(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// User tidak ditemukan
			return nil, errors.New("email atau password salah")
		}
		// Error database lainnya
		log.Printf("Error get user by email: %v", err)
		return nil, errors.New("terjadi kesalahan pada server")
	}

	// 2. Bandingkan password yang diberikan dengan hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Password tidak cocok
		return nil, errors.New("email atau password salah")
	}

	// 3. Jika password cocok, buatkan Token JWT
	tokenString, err := s.generateJWT(user.ID)
	if err != nil {
		log.Printf("Error generate JWT: %v", err)
		return nil, errors.New("gagal melakukan login")
	}

	// 4. Buat respons sukses
	userResponse := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	loginResponse := &models.LoginResponse{
		User:  userResponse,
		Token: tokenString,
	}

	return loginResponse, nil
}

// === Helper untuk Generate JWT ===
// generateJWT membuat token JWT baru untuk user
func (s *authService) generateJWT(userID string) (string, error) {
	// Tentukan 'claims' (data yang ingin kita simpan di token)
	// Kita juga set token ini kedaluwarsa dalam 7 hari
	claims := jwt.MapClaims{
		"sub": userID, // 'sub' (subject) adalah standar untuk ID user
		"iat": time.Now().Unix(), // 'iat' (issued at) waktu token dibuat
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 'exp' (expired)
	}

	// Buat token baru dengan claims dan metode signing HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token dengan secret key kita
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}