package repository

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
)

type PasswordResetStorage interface {
	CreatePasswordResetToken(ctx context.Context, userId string, passwordReset domain.PasswordReset) error
	DeleteAllTokens(ctx context.Context, userId string) error
	DeleteAllTokensOlderThan15min(ctx context.Context, userId string) error
	GetAmountTokensYoungerThan15Min(ctx context.Context, userId string) (int, error)
	CheckToken(ctx context.Context, tokenString string) error
	RevokeAllTokens(ctx context.Context, userId string) error
}
