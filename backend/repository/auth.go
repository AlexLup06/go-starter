package repository

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
)

type AuthStorage interface {
	CreateAuthProvider(ctx context.Context, authProvider domain.AuthProvider, userId string) error
	GetAuthProvider(ctx context.Context, userId string, method CreateUserMethod) (domain.AuthProvider, error)
	UpdateUserPassword(ctx context.Context, userId string, password string) error
}
