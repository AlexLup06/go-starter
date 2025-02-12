package service

import (
	"context"
	"fmt"

	"alexlupatsiy.com/personal-website/backend/domain"
	"alexlupatsiy.com/personal-website/backend/dto"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/repository"
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

	refreshToken, err, ttl := dto.GenerateJWT(userId, repository.REFRESH_TOKEN, s.JWTKey)
	if err != nil {
		return nil, -1, err
	}
	hashedToken := passwords.HashToken(refreshToken)

	jwt, err := dto.ParseJWT(refreshToken, s.JWTKey)
	if err != nil {
		return nil, -1, err
	}

	session := domain.Session{
		UserID:       userId,
		RefreshToken: hashedToken,
		IssuedAt:     jwt.Iat,
		ExpiresAt:    jwt.Exp,
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

func (s *SessionService) CreateAccessToken(ctx context.Context, userId string) (string, int64, error) {
	accessTokenString, err, ttl := dto.GenerateJWT(userId, repository.ACCESS_TOKEN, s.JWTKey)
	if err != nil {
		return "", -1, err
	}
	return accessTokenString, ttl, nil
}

func (s *SessionService) VerifyRefreshToken(ctx context.Context, refreshTokenString string) (bool, string, error) {
	accessToken, err := dto.ParseJWT(refreshTokenString, s.JWTKey)
	if err != nil {
		return false, "", err
	}

	// check against database

	return accessToken.IsExpired(), accessToken.Sub, nil
}

func (s *SessionService) VerfiyUserSession(ctx context.Context, userId string) error {
	session, err := s.sessionStorage.GetSessionByUserId(ctx, userId)
	if err != nil {
		return err
	}

	if session.IsExpired() {
		err := s.sessionStorage.RevokeSession(ctx, session.ID)
		if err != nil {
			return err
		}
		return fmt.Errorf("Refresh Token is expired. Need to login again")
	}
	return nil
}

func (s *SessionService) VerifyAccessToken(ctx context.Context, accessTokenString string) (bool, error) {
	accessToken, err := dto.ParseJWT(accessTokenString, s.JWTKey)
	if err != nil {
		return false, err
	}

	return accessToken.IsExpired(), nil
}

func (s *SessionService) GetSessionById(ctx context.Context, id string) (domain.Session, error) {
	session, err := s.sessionStorage.GetSessionById(ctx, id)
	if err != nil {
		return domain.Session{}, nil
	}
	return session, nil
}
