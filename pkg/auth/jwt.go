package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int64, secret string, expire int) (string, error) {
	// Generate a JWT token with userId and expiration time
	// Construct the claims(payload)
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Duration(expire) * time.Hour).Unix(),
	}
	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenStr string, secret string) (int64, error) {
	// Parse the token string and get userId
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	userId := int64(claims["userId"].(float64))

	return userId, nil
}
