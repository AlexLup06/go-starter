package service

import (
	"context"
	"fmt"
	"time"

	"alexlupatsiy.com/personal-website/backend/domain"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/repository"
	"github.com/golang-jwt/jwt"
)

type SessionService struct {
	sessionStorage repository.SessionStorage
	JWTKey         []byte
}

func NewSessionService(sessionStorage repository.SessionStorage, JWTKey []byte) *SessionService {
	return &SessionService{sessionStorage: sessionStorage, JWTKey: JWTKey}
}

// Only called in Login or Signup
func (s *SessionService) CreateRefreshToken(ctx context.Context, userId string) (*string, int64, error) {

	refreshToken, err, ttl := s.GenerateJWT(userId, repository.REFRESH_TOKEN)
	if err != nil {
		return nil, -1, err
	}
	hashedToken := passwords.HashToken(refreshToken)

	jwt, err := s.ParseJWT(refreshToken)
	if err != nil {
		return nil, -1, err
	}
	session := domain.Session{
		UserID:       userId,
		RefreshToken: hashedToken,
		IssuedAt:     time.Unix(int64(jwt["iat"].(float64)), 0),
		ExpiresAt:    time.Unix(int64(jwt["exp"].(float64)), 0),
		Revoked:      false,
		UserAgent:    "",
	}

	// revoke all sessions, even the active one, because we are getting a new session
	err = s.sessionStorage.RevokeAllSessions(ctx, userId)
	if err != nil {
		return nil, -1, err
	}

	err = s.sessionStorage.CreateSession(ctx, session)
	if err != nil {
		return nil, -1, err
	}

	return &refreshToken, ttl, nil
}

func (s *SessionService) CreateAccessToken(ctx context.Context, userId string) (*string, int64, error) {
	session, err := s.sessionStorage.GetSessionByUserId(ctx, userId)
	if err != nil {
		return nil, -1, err
	}

	if session.IsExpired() {
		err := s.sessionStorage.RevokeSession(ctx, session.ID)
		if err != nil {
			return nil, -1, err
		}
		return nil, -1, fmt.Errorf("Refresh Token is expired. Need to login again")
	}

	accessToken, err, ttl := s.GenerateJWT(userId, repository.ACCESS_TOKEN)
	if err != nil {
		return nil, -1, err
	}
	return &accessToken, ttl, nil
}

func (s *SessionService) VerifySession(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *SessionService) GetSessionById(ctx context.Context, id string) (domain.Session, error) {
	session, err := s.sessionStorage.GetSessionById(ctx, id)
	if err != nil {
		return domain.Session{}, nil
	}
	return session, nil
}

func (s *SessionService) GenerateJWT(userId string, tokenType repository.TokenType) (string, error, int64) {
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
	signedToken, err := token.SignedString(s.JWTKey)
	return signedToken, err, ttl
}

func (s *SessionService) ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
