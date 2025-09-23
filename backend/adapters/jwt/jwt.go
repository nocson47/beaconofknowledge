package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Get the secret from environment variable or fallback to a default (for development)
func JwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // Change this for production!
		//"p9$G7!kLz@wQ2#rT8sVbX1eF6uJmN4yZ"
	}
	return []byte(secret)
}

// Generate a JWT token with user_id claim
func GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		// 15 minutes expiry for access tokens
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString(JwtSecret())
}

// Parse and validate a JWT token, return user_id if valid
func ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JwtSecret(), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok {
			return int(userID), nil
		}
	}
	return 0, jwt.ErrInvalidKey
}
