package config

import (
	"fmt"
	"golang-test/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk akses database di file lain
var DB *gorm.DB

func ConnectDatabase() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using system default")
	}

	// Ambil value dari env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Susun DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Gagal koneksi ke database!")
	}

	err = database.AutoMigrate(&models.User{}, &models.Product{})
	if err != nil {
		// Ini bakal muncul di terminal kamu kalau migration gagal
		fmt.Println("Gagal Migrate:", err)
	} else {
		fmt.Println("Tabel User & Product Berhasil di-Migrate/Update!")
	}

	fmt.Println("Koneksi Database Berhasil!")

	DB = database
}
