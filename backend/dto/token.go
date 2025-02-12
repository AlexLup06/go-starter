package dto

import (
	"fmt"
	"time"

	"alexlupatsiy.com/personal-website/backend/repository"
	"github.com/golang-jwt/jwt"
)

type AccessToken struct {
	Sub string
	Exp time.Time
	Iat time.Time
	Iss string
}

func (a *AccessToken) IsExpired() bool {
	return a.Exp.Before(time.Now())
}

func (a *AccessToken) GetUserID() string {
	return a.Sub
}

func ParseJWT(tokenString string, JWTKey []byte) (AccessToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTKey, nil
	})

	if err != nil {
		return AccessToken{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessToken := AccessToken{
			Sub: claims["sub"].(string),
			Exp: time.Unix(int64(claims["exp"].(float64)), 0),
			Iat: time.Unix(int64(claims["iat"].(float64)), 0),
			Iss: claims["iss"].(string),
		}
		return accessToken, nil
	} else {
		return AccessToken{}, fmt.Errorf("invalid token")
	}
}

func GenerateJWT(userId string, tokenType repository.TokenType, JWTKey []byte) (string, error, int64) {
	var ttlMinutes time.Duration
	if tokenType.IsAccessToken() {
		ttlMinutes = 15
	} else {
		ttlMinutes = 20160 // 14 days
	}
	exp := time.Now().Add(ttlMinutes * time.Minute).Unix()
	iat := time.Now().Unix()
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": exp,
		"iat": iat,
		"iss": "app",
	}
	ttl := exp - iat

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWTKey)
	return signedToken, err, ttl
}
