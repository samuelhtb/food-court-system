package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk menampung koneksi database
var DB *gorm.DB

func ConnectDB() {
	// 1. Load file .env (hanya jika ada, berguna untuk development lokal)
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable bawaan sistem")
	}

	// 2. Ambil URL Database dari environment variable
	dsn := os.Getenv("SUPABASE_DB_URL")
	if dsn == "" {
		log.Fatal("Error: SUPABASE_DB_URL belum diatur di file .env")
	}

	// 3. Buka koneksi ke PostgreSQL (Supabase) menggunakan GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database Supabase: ", err)
	}

	log.Println("✅ Berhasil terhubung ke database Supabase!")

	// 4. Masukkan koneksi ke variabel global agar bisa dipakai oleh Repository nanti
	DB = database
}