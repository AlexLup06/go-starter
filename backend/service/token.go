package service

import "alexlupatsiy.com/personal-website/backend/helpers/token"

type TokenService struct {
	jwtKey []byte
}

func NewTokenService(jwtKey []byte) *TokenService {
	return &TokenService{jwtKey: jwtKey}
}

func (t *TokenService) GenerateUserInfoJWT(userInfo UserInfo, ttl int64) (string, *token.CustomClaims[UserInfo], int64, error) {
	tokenString, token, ttl, err := token.GenerateJWT(t.jwtKey, userInfo, 20160)
	if err != nil {
		return "", nil, -1, err
	}
	return tokenString, &token, ttl, err
}

func (t *TokenService) ParseUserInfoJWT(tokenString string) (*token.CustomClaims[UserInfo], error) {
	token, err := token.ParseJWT[UserInfo](t.jwtKey, tokenString)
	if err != nil {
		return nil, err
	}
	return token, err
}
