package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nama-github-kamu/bookwise-backend/internal/handlers" // Kita akan pinjam response helper
)

// Tipe custom untuk context key
type contextKey string
const userContextKey = contextKey("userID")

// AuthMiddleware adalah penjaga endpoint kita
func AuthMiddleware(next http.Handler) http.Handler {
	// Ambil JWT Secret dari .env
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("FATAL: JWT_SECRET tidak diset di .env")
	}
	
	// 'fn' adalah fungsi yang akan dieksekusi setiap ada request
	fn := func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Header Authorization kosong")
			return
		}

		// 2. Cek format "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Format token tidak valid")
			return
		}

		// 3. Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing-nya benar
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler // Error aneh
			}
			return []byte(jwtSecret), nil
		})

		// 4. Cek error atau token tidak valid
		if err != nil || !token.Valid {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Token tidak valid atau kedaluwarsa")
			return
		}

		// 5. Ambil 'claims' (data di dalam token)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Token claims tidak valid")
			return
		}

		// 6. Ambil UserID ('sub') dari claims
		userID, ok := claims["sub"].(string)
		if !ok {
			handlers.RespondWithError(w, http.StatusUnauthorized, "UserID tidak ditemukan di token")
			return
		}

		// 7. Sukses! Simpan UserID di context
		// Ini berguna jika ada endpoint yang butuh tahu "siapa" yang request
		ctx := context.WithValue(r.Context(), userContextKey, userID)
		
		// 8. Lanjutkan ke handler berikutnya (misal: CreateCategory)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// GetUserIDFromContext adalah helper untuk mengambil UserID dari context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userContextKey).(string)
	return userID, ok
}