package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL
)

// InitDB membuat dan mengembalikan koneksi database pool
func InitDB() *sql.DB {
	// Baca DSN (Data Source Name) dari .env
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("FATAL: DATABASE_URL tidak diset di .env")
	}

	// Buka koneksi ke MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("FATAL: Tidak bisa membuka koneksi database: %v", err)
	}

	// Coba ping ke database untuk memastikan koneksi valid
	if err := db.Ping(); err != nil {
		log.Fatalf("FATAL: Tidak bisa terhubung ke database: %v", err)
	}

	log.Println("Koneksi database berhasil.")
	return db
}