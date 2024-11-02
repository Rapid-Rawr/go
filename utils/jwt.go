package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Make JwtKey an exported variable by renaming it
var JwtKey = []byte("secret_key")

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(JwtKey)
}
