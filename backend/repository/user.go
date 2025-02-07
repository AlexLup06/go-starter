package repository

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
)

type UpdateUser struct {
	username *string
	password *string
}

type UserStorage interface {
	CreateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, request UpdateUser) error
}
