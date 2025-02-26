package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims[T any] struct {
	Extra T `json:"extra"`
	jwt.RegisteredClaims
}

func GenerateJWT[T any](secretKey []byte, customPayload T, ttlMinutes time.Duration) (string, CustomClaims[T], int64, error) {
	exp := time.Now().Add(ttlMinutes * time.Minute)

	claims := CustomClaims[T]{
		Extra: customPayload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", CustomClaims[T]{}, 0, err
	}

	return signedToken, claims, exp.Unix(), nil
}

func ParseJWT[T any](secretKey []byte, tokenString string) (*CustomClaims[T], error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims[T]{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims[T])
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
