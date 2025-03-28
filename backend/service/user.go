package service

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
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

type SignInWithGoogle struct {
	IdToken   string `form:"credential"`
	CSRFToken string `form:"g_csrf_token"`
}

func (u *UserService) CreateUser(ctx context.Context, email, name string) (UserInfo, error) {
	user := domain.User{
		Username: name,
		Email:    &email,
	}

	createdUser, err := u.userStorage.CreateUser(ctx, user)
	if err != nil {
		return UserInfo{}, err
	}

	userInfo := UserInfo{
		UserId:   createdUser.ID,
		Username: createdUser.Username,
		Email:    *createdUser.Email,
	}

	return userInfo, nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserService) UpdateUserEmail(ctx context.Context, userId, email string) error {
	err := u.userStorage.UpdateUserEmail(ctx, userId, email)
	if err != nil {
		return err
	}
	return nil
}
