package db

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type userDb struct{}

func NewUserDb() repository.UserStorage {
	return &userDb{}
}

func (c *userDb) CreateUser(ctx context.Context, user domain.User) error {
	return nil
}

func (c *userDb) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (c *userDb) UpdateUser(ctx context.Context, request repository.UpdateUser) error {
	return nil
}
