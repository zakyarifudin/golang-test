package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Gin
	r := gin.Default()

	// 2. Bikin Route sederhana
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "sukses",
			"message": "Halo Zaky! Backend Go kamu sudah jalan di M2 bre",
		})
	})

	// 3. Jalankan server di port 8001
	r.Run(":8001")
}
