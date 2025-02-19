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
	jwtKey         []byte
}

func NewSessionService(sessionStorage repository.SessionStorage, jwtKey []byte) *SessionService {
	return &SessionService{sessionStorage: sessionStorage, jwtKey: jwtKey}
}

// Only called in Login or Signup
func (s *SessionService) CreateRefreshToken(ctx context.Context, userInfo UserInfo) (string, int64, error) {

	refreshToken, err, ttl := s.GenerateJWT(userInfo, repository.REFRESH_TOKEN.Type)
	if err != nil {
		return "", -1, err
	}
	hashedToken := passwords.HashToken(refreshToken)

	jwt, err := s.ParseJWT(refreshToken)
	if err != nil {
		return "", -1, err
	}

	session := domain.Session{
		UserID:       userInfo.UserId,
		RefreshToken: hashedToken,
		IssuedAt:     jwt.Iat,
		ExpiresAt:    jwt.Exp,
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

	return refreshToken, ttl, nil
}

func (s *SessionService) CreateAccessToken(ctx context.Context, userInfo UserInfo) (string, int64, error) {
	accessTokenString, err, ttl := s.GenerateJWT(userInfo, repository.ACCESS_TOKEN.Type)
	if err != nil {
		return "", -1, err
	}
	return accessTokenString, ttl, nil
}

func (s *SessionService) VerifyRefreshToken(ctx context.Context, refreshTokenString string) (bool, UserInfo, error) {
	refreshToken, err := s.ParseJWT(refreshTokenString)
	if err != nil {
		return false, UserInfo{}, err
	}
	hashedRefreshToken := passwords.HashToken(refreshTokenString)

	isValid, err := s.sessionStorage.ValidateRefreshToken(ctx, hashedRefreshToken, refreshToken.UserID)

	userInfo := UserInfo{
		UserId:   refreshToken.UserID,
		Username: refreshToken.Username,
		Email:    refreshToken.Email,
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
	accessToken, err := s.ParseJWT(accessTokenString)
	if err != nil {
		return false, "", err
	}

	return accessToken.IsExpired(), accessToken.GetUserID(), nil
}

func (s *SessionService) GetSessionById(ctx context.Context, id string) (domain.Session, error) {
	session, err := s.sessionStorage.GetSessionById(ctx, id)
	if err != nil {
		return domain.Session{}, nil
	}
	return session, nil
}

func (s SessionService) RevokeAllSessions(ctx context.Context, userId string) error {
	// revoke all sessions, even the active one, because we are getting a new session
	err := s.sessionStorage.RevokeAllSessions(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

type UserInfo struct {
	UserId   string
	Username string
	Email    string
}

type Token struct {
	UserID   string
	Username string
	Email    string
	Exp      time.Time
	Iat      time.Time
	Iss      string
}

// IsExpired checks if the token is expired
func (a *Token) IsExpired() bool {
	return a.Exp.Before(time.Now())
}

// GetUserID retrieves the user ID
func (a *Token) GetUserID() string {
	return a.UserID
}

// ParseJWT decodes and verifies the JWT
// TODO: error handling. If fields dont exist
func (s *SessionService) ParseJWT(tokenString string) (Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
	})

	if err != nil {
		return Token{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		token := Token{
			UserID:   claims["user_id"].(string),
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			Exp:      time.Unix(int64(claims["exp"].(float64)), 0),
			Iat:      time.Unix(int64(claims["iat"].(float64)), 0),
			Iss:      claims["iss"].(string),
		}
		return token, nil
	} else {
		return Token{}, fmt.Errorf("invalid token")
	}
}

// GenerateJWT creates a JWT with user details
func (s *SessionService) GenerateJWT(userInfo UserInfo, tokenType string) (string, error, int64) {
	var ttlMinutes time.Duration
	if tokenType == repository.ACCESS_TOKEN.Type {
		ttlMinutes = 15
	} else {
		ttlMinutes = 20160 // 14 days
	}
	exp := time.Now().Add(ttlMinutes * time.Minute).Unix()
	iat := time.Now().Unix()
	claims := jwt.MapClaims{
		"user_id":  userInfo.UserId,
		"username": userInfo.Username,
		"email":    userInfo.Email,
		"exp":      exp,
		"iat":      iat,
		"iss":      "app",
	}
	ttl := exp - iat

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtKey)
	return signedToken, err, ttl
}
