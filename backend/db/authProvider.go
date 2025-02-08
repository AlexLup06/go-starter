package db

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type authDb struct{}

func NewAuthDb() repository.AuthStorage {
	return &authDb{}
}

func (a *authDb) CreateAuthProvider(ctx context.Context, authProvider domain.AuthProvider, userId string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}
	var authProviders []domain.AuthProvider
	err = db.Model(domain.AuthProvider{}).Where("user_id = ?", userId).Find(&authProviders).Error

	if err == nil {
		for _, authProvider := range authProviders {
			// TODO: make geneal
			if authProvider.Method == repository.METHOD_EMAIL.Method {
				return customErrors.ErrEmailAndProviderExist
			}
		}
	}

	err = db.Model(domain.AuthProvider{}).Create(&authProvider).Error

	if err != nil {
		if customErrors.IsUniqueConstraintViolationError(err) {
			return customErrors.ErrEmailAndProviderExist
		}
		return err
	}

	return nil
}

func (a *authDb) GetAuthProvider(ctx context.Context, userId string, method repository.CreateUserMethod) (domain.AuthProvider, error) {
	db, err := getContextDb(ctx)
	if err != nil {
		return domain.AuthProvider{}, err
	}
	var authProvider domain.AuthProvider
	err = db.Model(domain.AuthProvider{}).Where("user_id = ?", userId).Where("method = ?", method.Method).First(&authProvider).Error

	if err != nil {
		return domain.AuthProvider{}, customErrors.ErrAuthProviderDoesNotExist
	}

	return authProvider, nil
}
