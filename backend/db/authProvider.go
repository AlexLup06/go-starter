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

func (a *authDb) CreateAuthProvider(ctx context.Context, authProvider domain.AuthProvider) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}
	var usersAuthProviders []domain.AuthProvider
	err = db.Model(domain.AuthProvider{}).Where("user_id = ?", authProvider.UserID).Find(&usersAuthProviders).Error
	if err == nil {
		for _, uasersAP := range usersAuthProviders {
			if authProvider.Method == uasersAP.Method {
				return customErrors.ErrUserWithAuthProviderExist
			}
		}
	}

	err = db.Model(domain.AuthProvider{}).Create(&authProvider).Error

	if err != nil {
		return err
	}

	return nil
}

func (a *authDb) GetAuthProviderByUserId(ctx context.Context, userId string, method repository.CreateUserMethod) (domain.AuthProvider, error) {
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

func (a *authDb) GetAuthProviderByProviderId(ctx context.Context, providerId string, method repository.CreateUserMethod) (domain.AuthProvider, error) {
	db, err := getContextDb(ctx)
	if err != nil {
		return domain.AuthProvider{}, err
	}
	var authProvider domain.AuthProvider
	err = db.Model(domain.AuthProvider{}).Where("provider_user_id = ?", providerId).Where("method = ?", method.Method).First(&authProvider).Error

	if err != nil {
		return domain.AuthProvider{}, customErrors.ErrAuthProviderDoesNotExist
	}

	return authProvider, nil
}

func (a *authDb) UpdateUserPassword(ctx context.Context, userId string, hashedPassword string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	err = db.Model(&domain.AuthProvider{}).Where("user_id = ?", userId).Update("password_hash", hashedPassword).Error
	if err != nil {
		return err
	}

	return nil
}
