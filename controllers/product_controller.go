package controllers

import (
	"golang-test/config"
	"golang-test/models"
	"net/http"
	"strconv"

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
		Barcode string `json:"barcode" binding:"required"` // Tambahkan ini
		Name    string `json:"name" binding:"required"`
		Price   int    `json:"price" binding:"required"`
		Stock   int    `json:"stock" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Barcode: input.Barcode, // Masukkan ke struct
		Name:    input.Name,
		Price:   input.Price,
		Stock:   input.Stock,
	}

	// Simpan ke DB, GORM otomatis cek kalau barcode-nya duplikat bakal error
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan! Barcode mungkin sudah terdaftar.", "error_detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang berhasil ditambah!", "data": product})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang nggak ada!"})
		return
	}

	// Bind data baru
	c.ShouldBindJSON(&product)
	config.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{"message": "Data diupdate!", "data": product})
}

// Hapus Barang
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Barang terhapus!"})
}
