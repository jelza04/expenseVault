package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtSecret is the signing key for JWT tokens.
var JwtSecret []byte

// GenerateToken creates a JWT for a username.
func GenerateToken(username string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
