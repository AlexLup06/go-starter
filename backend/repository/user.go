package repository

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
)

// type UpdateUser struct {
// 	username *string
// 	password *string
// }

type CreateUserMethod struct {
	Method string
}

var (
	METHOD_EMAIL  = CreateUserMethod{Method: "email"}
	METHOD_APPLE  = CreateUserMethod{Method: "apple"}
	METHOD_GOOGLE = CreateUserMethod{Method: "google"}
)

type UserStorage interface {
	CreateUser(ctx context.Context, user domain.User, authProvider domain.AuthProvider) (domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
