package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint
	Role   string
	jwt.RegisteredClaims
}

// Ambil key dari env
func GetJwtKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("default_secret_key") // Fallback kalau env kosong
	}
	return []byte(secret)
}

func GenerateToken(userID uint, role string) (string, error) {
	// Ubah jadi 30 Hari
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJwtKey())
}

func GenerateRefreshToken(userID uint) (string, error) {
	// Refresh token harus lebih lama dari access token, misal 60 hari
	expirationTime := time.Now().Add(60 * 24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJwtKey())
}
