package controllers

import (
	"golang-test/config"
	"golang-test/helpers"
	"golang-test/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password biar aman (Gak boleh simpan plain text!)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	config.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dibuat!"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah!"})
		return
	}

	var user models.User
	// Cari user
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/Password salah!"})
		return
	}

	// Cek password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/Password salah!"})
		return
	}

	// Generate Access Token (Sekarang 30 Hari sesuai helper tadi)
	accessToken, err := helpers.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal buat access token"})
		return
	}

	// Generate Refresh Token (60 Hari)
	refreshToken, _ := helpers.GenerateRefreshToken(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login berhasil, token aktif selama 30 hari!",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    "30 Days", // Tambahin info ini biar Frontend tau
	})
}

func Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// 1. Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token mana, Bos?"})
		return
	}

	// 2. Validasi & Parsing Refresh Token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return helpers.GetJwtKey(), nil
	})

	// Cek apakah token valid atau sudah basi (expired)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token palsu atau udah basi!"})
		return
	}

	// 3. Ambil data User ID dari Token (disimpan di Claims 'sub')
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal baca isi token"})
		return
	}

	userIDStr := claims["sub"].(string)

	// 4. Cari User di DB buat mastiin usernya masih aktif & ambil Role-nya
	var user models.User
	if err := config.DB.First(&user, userIDStr).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User udah gak ada di database!"})
		return
	}

	// 5. Generate Access Token baru (30 Hari)
	newAccessToken, _ := helpers.GenerateToken(user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Token berhasil diperpanjang!",
		"access_token": newAccessToken,
	})
}
