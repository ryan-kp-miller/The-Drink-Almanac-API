package service

import (
	"time"

	"the-drink-almanac-api/errors"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	// CreateNewToken generates a new token with the provided userId stored in it and an expiry based on the provided number of minutes
	CreateNewToken(userId string, ttlMinutes int) (string, error)

	// ValidateToken verifies that the token is valid and returns the userId if it's valid
	ValidateToken(string) (string, error)
}

type JwtAuthService struct {
	authSecretKey []byte
}

func (s JwtAuthService) CreateNewToken(userId string, expiryDurationMinutes int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    time.Now().Add(time.Duration(expiryDurationMinutes) * time.Minute).Unix(),
		"userId": userId,
	})

	tokenString, err := token.SignedString(s.authSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s JwtAuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.NewInvalidAuthTokenError("invalid token format")
		}
		return s.authSecretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId := claims["userId"].(string)
		return userId, nil
	} else {
		return "", errors.NewInvalidAuthTokenError("token is no longer valid")
	}
}

func NewJwtAuthService(secretKey string) JwtAuthService {
	return JwtAuthService{
		authSecretKey: []byte(secretKey),
	}
}
