package util

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims jwt.Claims) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (jwt.Claims, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
