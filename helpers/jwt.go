package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key - Ganti sama string random kamu sendiri!
var JWT_KEY = []byte("rahasia_zaky_jwt")

type JWTClaims struct {
	UserID uint
	Role   string
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
	// Token berlaku 24 jam
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_KEY)
}
