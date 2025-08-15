package services

import (
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("mysecretkey")

func ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}

	userID := int(claims["user_id"].(float64))
	return userID, nil
}