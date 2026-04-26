package controllers

import (
	"encoding/base64"
	"fmt"
	"golang-test/config"
	"golang-test/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []models.Product
	var totalCount int64

	// Ganti name jadi search
	searchQuery := c.Query("search")
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Product{})

	// Logika pencarian yang lebih fleksibel
	if searchQuery != "" {
		// Sekarang bisa cari berdasarkan Nama ATAU Barcode di satu kolom search
		query = query.Where("name LIKE ? OR barcode LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	query.Count(&totalCount)

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal tarik data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Data berhasil ditarik",
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
		"data":        products,
	})
}

func CreateProduct(c *gin.Context) {
	var input struct {
		Barcode string `json:"barcode" binding:"required"`
		Name    string `json:"name" binding:"required"`
		Price   int    `json:"price" binding:"required"`
		Stock   int    `json:"stock" binding:"required"`
		Image   string `json:"image"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filePath := ""
	if input.Image != "" {
		path, err := saveBase64Image(input.Image, input.Barcode)
		if err == nil {
			filePath = path
		}
	}

	product := models.Product{
		Barcode: input.Barcode,
		Name:    input.Name,
		Price:   input.Price,
		Stock:   input.Stock,
		Image:   filePath,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Barcode sudah terdaftar!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dibuat", "data": product})
}

func GetProductByBarcode(c *gin.Context) {
	barcode := c.Param("barcode") // Ambil dari URL
	var product models.Product

	if err := config.DB.Where("barcode = ?", barcode).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan!", "error_detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// Update Barang
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var input struct {
		Barcode string `json:"barcode"`
		Name    string `json:"name"`
		Price   int    `json:"price"`
		Stock   int    `json:"stock"`
		Image   string `json:"image"`
	}

	c.ShouldBindJSON(&input)

	// Logic Image: Jika input.Image berisi Base64 (bukan path uploads/)
	if input.Image != "" && !strings.HasPrefix(input.Image, "uploads/") {
		// Hapus foto lama jika ada
		if product.Image != "" {
			os.Remove(product.Image)
		}

		// Simpan foto baru
		newPath, err := saveBase64Image(input.Image, product.Barcode)
		if err == nil {
			product.Image = newPath
		}
	}

	// Update fields
	if input.Barcode != "" {
		product.Barcode = input.Barcode
	}
	if input.Name != "" {
		product.Name = input.Name
	}
	if input.Price != 0 {
		product.Price = input.Price
	}
	if input.Stock != 0 {
		product.Stock = input.Stock
	}

	config.DB.Save(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Update berhasil", "data": product})
}

// Hapus Barang
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// Cari dulu buat dapet path fotonya
	if err := config.DB.First(&product, id).Error; err == nil {
		// Hapus file fotonya dari storage MacBook
		if product.Image != "" {
			os.Remove(product.Image)
		}
	}

	// Hapus data dari DB
	config.DB.Delete(&models.Product{}, id)

	c.JSON(http.StatusOK, gin.H{"message": "Barang dan fotonya sudah musnah!"})
}

// Helper untuk simpan Base64 ke File
func saveBase64Image(base64Str string, barcode string) (string, error) {
	// 1. Cek prefix data:image/...
	i := strings.Index(base64Str, ",")
	if i == -1 {
		return "", fmt.Errorf("invalid base64 format")
	}

	// 2. Decode stringnya
	rawDecodedText, err := base64.StdEncoding.DecodeString(base64Str[i+1:])
	if err != nil {
		return "", err
	}

	// 3. Buat nama file unik
	fileName := fmt.Sprintf("uploads/%s_%d.png", barcode, time.Now().Unix())

	// 4. Tulis ke storage
	err = os.WriteFile(fileName, rawDecodedText, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
