package service

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type LoginWithEmailRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type AuthService struct {
	authStorage    repository.AuthStorage
	userService    *UserService
	sessionService *SessionService
}

func NewAuthService(authStorage repository.AuthStorage, userService *UserService, sessionService *SessionService) *AuthService {
	return &AuthService{authStorage: authStorage, userService: userService, sessionService: sessionService}
}

func (a *AuthService) LoginWithEmail(ctx context.Context, request LoginWithEmailRequest) (domain.User, error) {
	user, err := a.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return domain.User{}, err
	}

	emailAuth, err := a.authStorage.GetAuthProvider(ctx, user.ID, repository.METHOD_EMAIL)
	if err != nil {
		// the auth provider "email" does not exist
		return domain.User{}, err
	}

	if !passwords.IsSamePassword(request.Password, *emailAuth.PasswordHash) {
		return domain.User{}, customErrors.NewUnauthorizedError("invalid password")
	}

	return user, nil
}
