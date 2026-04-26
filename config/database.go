package config

import (
	"fmt"
	"golang-test/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk akses database di file lain
var DB *gorm.DB

func ConnectDatabase() {
	// Format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	// Di DBngin biasanya user: root, password: kosong, port: 3306
	dsn := "root:@tcp(127.0.0.1:3306)/golang_test?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Gagal koneksi ke database! Pastikan DBngin sudah jalan.")
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
