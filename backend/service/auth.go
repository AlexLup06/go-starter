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

type RequestPasswordResetRequest struct {
	Email string `form:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Password string `form:"password" binding:"required"`
	Token    string `form:"token" binding:"required"`
}

type AuthService struct {
	authStorage  repository.AuthStorage
	userService  *UserService
	tokenService *TokenService
}

func NewAuthService(authStorage repository.AuthStorage, userService *UserService, tokenService *TokenService) *AuthService {
	return &AuthService{authStorage: authStorage, userService: userService, tokenService: tokenService}
}

func (a *AuthService) LoginWithEmail(ctx context.Context, request LoginWithEmailRequest) (UserInfo, error) {
	user, err := a.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return UserInfo{}, err
	}

	emailAuth, err := a.authStorage.GetAuthProvider(ctx, user.ID, repository.METHOD_EMAIL)
	if err != nil {
		// the auth provider "email" does not exist
		return UserInfo{}, err
	}

	if !passwords.IsSamePassword(request.Password, *emailAuth.PasswordHash) {
		return UserInfo{}, customErrors.NewUnauthorizedError("invalid password")
	}

	userInfo := UserInfo{
		UserId:   user.ID,
		Username: user.Username,
		Email:    *user.Email,
	}

	return userInfo, nil
}

func (a *AuthService) UpdateUserPassword(ctx context.Context, userId string, resetPasswordRequest ResetPasswordRequest) error {
	hashedPassword, err := passwords.HashPassword(resetPasswordRequest.Password)
	if err != nil {
		return err
	}
	err = a.authStorage.UpdateUserPassword(ctx, userId, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}
