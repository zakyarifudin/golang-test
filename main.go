package main

import (
	"golang-test/config"
	"golang-test/controllers"
	"golang-test/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Jalankan koneksi database sebelum server start
	config.ConnectDatabase()
	// 1. Inisialisasi Gin
	r := gin.Default()

	// 2. Bikin Route sederhana
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "sukses",
			"message": "Halo Zaky! Backend Go kamu sudah jalan di M2 bre",
		})
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.Refresh)
	// Protected Routes (Harus pake Token)
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/products", controllers.GetProducts)
		protected.POST("/products", controllers.CreateProduct)
		protected.GET("/products/barcode/:barcode", controllers.GetProductByBarcode)
		protected.PUT("/products/:id", controllers.UpdateProduct)
		protected.DELETE("/products/:id", controllers.DeleteProduct)
	}

	// 3. Jalankan server di port 8001
	r.Run(":8001")
}
