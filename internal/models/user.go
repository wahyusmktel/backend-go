package models

import "time"

// RegisterRequest adalah model untuk payload JSON saat registrasi
// Sesuai dengan form 'handleRegister'
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest adalah model untuk payload JSON saat login
// Sesuai dengan form 'handleLogin'
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// User adalah model yang merepresentasikan data di tabel 'users'
// Ini untuk penggunaan internal (repository, service)
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Tanda '-' SANGAT PENTING!
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
// Penjelasan `json:"-"`:
// Ini memberitahu Go untuk SELALU mengabaikan field 'Password'
// saat mengubah struct ini menjadi JSON.
// Ini adalah langkah keamanan agar hash password tidak pernah
// bocor ke frontend.

// UserResponse adalah model yang kita kirim kembali ke klien (tanpa password)
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginResponse adalah model yang dikirim setelah login sukses
// Frontend akan menerima ini
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"` // Token JWT akan kita buat nanti
}