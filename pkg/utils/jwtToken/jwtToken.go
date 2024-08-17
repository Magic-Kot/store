package jwtToken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	signingKey = "hs1dfjKhl0iLLLhfjH7"
	tokenTTL   = 2 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	ID int `json:"id"`
}

// GenerateToken - генерация токена
func GenerateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		ID: id,
	})

	return token.SignedString([]byte(signingKey))
}

// ParseToken - парсинг токена
func ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.ID, nil
}
