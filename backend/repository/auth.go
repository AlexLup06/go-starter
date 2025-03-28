package repository

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
)

type AuthStorage interface {
	CreateAuthProvider(ctx context.Context, authProvider domain.AuthProvider) error
	GetAuthProviderByUserId(ctx context.Context, userId string, method CreateUserMethod) (domain.AuthProvider, error)
	GetAuthProviderByProviderId(ctx context.Context, providerId string, method CreateUserMethod) (domain.AuthProvider, error)
	UpdateUserPassword(ctx context.Context, userId string, password string) error
}
