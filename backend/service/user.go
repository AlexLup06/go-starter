package service

import (
	"context"
	"fmt"

	"alexlupatsiy.com/personal-website/backend/domain"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/helpers/passwords"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type UserService struct {
	userStorage repository.UserStorage
	authStorage repository.AuthStorage
}

func NewUserService(userDb repository.UserStorage, authStorage repository.AuthStorage) *UserService {
	return &UserService{userStorage: userDb, authStorage: authStorage}
}

type SignUpWithEmailRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (u *UserService) CreateUserWithEmail(ctx context.Context, request SignUpWithEmailRequest) error {
	hashedPassword, err := passwords.HashPassword(request.Password)
	if err != nil {
		return fmt.Errorf("error hashing the user's password: %w", err)
	}

	emailAuthProvider := domain.AuthProvider{
		Method:       repository.METHOD_EMAIL.Method,
		PasswordHash: &hashedPassword,
	}
	user := domain.User{
		Username: "Alex", // TODO: Generate Random Username
		Email:    &request.Email,
		AuthProviders: []domain.AuthProvider{
			emailAuthProvider,
		},
	}

	existingUser, err := u.userStorage.CreateUser(ctx, user, emailAuthProvider)
	if err == nil {
		return nil
	}

	// user signed up with social login and now wants to create password login
	if err == customErrors.ErrEmailExists {
		err = u.authStorage.CreateAuthProvider(ctx, emailAuthProvider, existingUser.ID)
		if err != nil {
			// the auth provider already exists for the email
			return err
		}
		return nil
	}

	if err == customErrors.ErrEmailAndProviderExist {
		return customErrors.NewValidationError(fmt.Sprintf("user %q already exists with the provider", *user.Email))
	}

	return err
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
