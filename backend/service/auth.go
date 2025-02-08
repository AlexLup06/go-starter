package service

import (
	"context"

	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type LoginWithEmailRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type AuthService struct {
	authStorage repository.AuthStorage
	userService UserService
}

func NewAuthService(authStorage repository.AuthStorage, userService *UserService) *AuthService {
	return &AuthService{authStorage: authStorage, userService: *userService}
}

func (a *AuthService) LoginWithEmail(ctx context.Context, request LoginWithEmailRequest) error {
	user, err := a.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	emailAuth, err := a.authStorage.GetAuthProvider(ctx, user.ID, repository.METHOD_EMAIL)
	if err != nil {
		// the auth provider "email" does not exist
		return err
	}

	if !passwords.IsSamePassword(request.Password, *emailAuth.PasswordHash) {
		return customErrors.NewUnauthorizedError("invalid password")
	}

	return nil
}
