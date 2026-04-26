package middlewares

import (
	"golang-test/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Butuh token buat masuk, Bos!"})
			c.Abort()
			return
		}

		// 2. Formatnya harus "Bearer <token>"
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. Parsing & Validasi Token
		token, err := jwt.ParseWithClaims(tokenString, &helpers.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return helpers.GetJwtKey(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token palsu atau udah basi!"})
			c.Abort()
			return
		}

		// 4. Kalau oke, simpan data user ke context biar bisa dipake di controller
		claims, _ := token.Claims.(*helpers.JWTClaims)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
