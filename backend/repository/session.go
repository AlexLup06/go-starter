package repository

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
)

type SessionStorage interface {
	CreateSession(ctx context.Context, session domain.Session) error
	DeleteSession(ctx context.Context, sessionId string) error
	GetSessionById(ctx context.Context, sessionId string) (domain.Session, error)
	GetSessionByUserId(ctx context.Context, sessionId string) (domain.Session, error)
	RevokeSession(ctx context.Context, sessionId string) error
	RevokeAllSessions(ctx context.Context, userId string) error
}

type CookieType struct {
	Type string
}

func (t *CookieType) IsRefreshToken() bool {
	return t.Type == "refresh_token"
}
func (t *CookieType) IsAccessToken() bool {
	return t.Type == "access_token"
}

type TokenType struct {
	Type string
}

func (t *TokenType) IsRefreshToken() bool {
	return t.Type == "refresh_token"
}
func (t *TokenType) IsAccessToken() bool {
	return t.Type == "access_token"
}

var (
	REFRESH_TOKEN = TokenType{Type: "refresh_token"}
	ACCESS_TOKEN  = TokenType{Type: "access_token"}

	REFRESH_COOKIE = CookieType{Type: "refresh_cookie"}
	ACCESS_COOKIE  = CookieType{Type: "access_cookie"}
)
