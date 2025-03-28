package db

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type userDb struct{}

func NewUserDb() repository.UserStorage {
	return &userDb{}
}

func (u *userDb) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	db, err := getContextDb(ctx)
	if err != nil {
		return domain.User{}, err
	}

	var existingUser domain.User
	err = db.Model(domain.User{}).Where("LOWER(email) = LOWER(?)", user.Email).First(&existingUser).Error

	if err == nil {
		return existingUser, customErrors.ErrEmailExists
	}

	err = db.Model(domain.User{}).Create(&user).Error

	if err != nil {
		if customErrors.IsUniqueConstraintViolationError(err) {
			return domain.User{}, customErrors.ErrUserWithAuthProviderExist
		}
		return domain.User{}, err
	}

	return user, nil
}

func (u *userDb) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	db, err := getContextDb(ctx)
	if err != nil {
		return domain.User{}, err
	}

	var existingUser domain.User
	err = db.Model(domain.User{}).Where("LOWER(email) = LOWER(?)", email).First(&existingUser).Error

	if err != nil {
		return existingUser, customErrors.ErrUserDoesNotExist
	}
	return existingUser, nil
}

func (u *userDb) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (u *userDb) UpdateUserEmail(ctx context.Context, userId, email string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}
	err = db.Model(&domain.User{}).Where("id = ?", userId).Update("email", email).Error
	if err != nil {
		return err
	}

	return nil
}
