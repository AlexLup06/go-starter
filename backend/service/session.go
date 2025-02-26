package service

import (
	"context"
	"fmt"
	"time"

	"alexlupatsiy.com/personal-website/backend/domain"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type UserInfo struct {
	UserId   string
	Username string
	Email    string
}

type SessionService struct {
	sessionStorage repository.SessionStorage
	tokenService   *TokenService
}

func NewSessionService(sessionStorage repository.SessionStorage, tokenService *TokenService) *SessionService {
	return &SessionService{sessionStorage: sessionStorage, tokenService: tokenService}
}

// Only called in Login or Signup
func (s *SessionService) CreateRefreshToken(ctx context.Context, userInfo UserInfo) (string, int64, error) {

	refreshTokenString, refreshToken, ttl, err := s.tokenService.GenerateUserInfoJWT(userInfo, 20160)
	if err != nil {
		return "", -1, err
	}

	hashedToken := passwords.HashToken(refreshTokenString)

	session := domain.Session{
		UserID:       userInfo.UserId,
		RefreshToken: hashedToken,
		IssuedAt:     refreshToken.IssuedAt.Time,
		ExpiresAt:    refreshToken.ExpiresAt.Time,
		Revoked:      false,
		UserAgent:    "",
	}

	// revoke all sessions, even the active one, because we are getting a new session
	err = s.RevokeAllSessions(ctx, userInfo.UserId)
	if err != nil {
		return "", -1, err
	}

	err = s.sessionStorage.CreateSession(ctx, session)
	if err != nil {
		return "", -1, err
	}

	return refreshTokenString, ttl, nil
}

func (s *SessionService) CreateAccessToken(ctx context.Context, userInfo UserInfo) (string, int64, error) {
	accessTokenString, _, ttl, err := s.tokenService.GenerateUserInfoJWT(userInfo, 15)
	if err != nil {
		return "", -1, err
	}
	return accessTokenString, ttl, nil
}

func (s *SessionService) VerifyRefreshToken(ctx context.Context, refreshTokenString string) (bool, UserInfo, error) {
	refreshToken, err := s.tokenService.ParseUserInfoJWT(refreshTokenString)
	if err != nil {
		return false, UserInfo{}, err
	}
	hashedRefreshToken := passwords.HashToken(refreshTokenString)

	isValid, err := s.sessionStorage.ValidateRefreshToken(ctx, hashedRefreshToken, refreshToken.Extra.UserId)

	userInfo := UserInfo{
		UserId:   refreshToken.Extra.UserId,
		Username: refreshToken.Extra.Username,
		Email:    refreshToken.Extra.Email,
	}

	return isValid, userInfo, nil
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

func (s *SessionService) VerifyAccessToken(ctx context.Context, accessTokenString string) (bool, string, error) {
	accessToken, err := s.tokenService.ParseUserInfoJWT(accessTokenString)
	if err != nil {
		return false, "", err
	}

	isValid := accessToken.VerifyExpiresAt(time.Now(), true)

	return isValid, accessToken.Extra.UserId, nil
}

func (s *SessionService) GetSessionById(ctx context.Context, id string) (domain.Session, error) {
	session, err := s.sessionStorage.GetSessionById(ctx, id)
	if err != nil {
		return domain.Session{}, nil
	}
	return session, nil
}

func (s *SessionService) RevokeAllSessions(ctx context.Context, userId string) error {
	// revoke all sessions, even the active one, because we are getting a new session
	err := s.sessionStorage.RevokeAllSessions(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionService) ParseUserInfo(refreshTokenString string) (UserInfo, error) {
	token, err := s.tokenService.ParseUserInfoJWT(refreshTokenString)
	return token.Extra, err
}
