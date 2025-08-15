package services

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(secret),
	}
}

func (a *AuthService) ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
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