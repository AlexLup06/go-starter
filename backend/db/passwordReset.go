package db

import (
	"context"
	"fmt"
	"time"

	"alexlupatsiy.com/personal-website/backend/domain"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type passwordResetDb struct{}

func NewPasswordResetDb() repository.PasswordResetStorage {
	return &passwordResetDb{}
}

func (p *passwordResetDb) CreatePasswordResetToken(ctx context.Context, userId string, passwordReset domain.PasswordReset) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	err = db.Model(&domain.PasswordReset{}).Create(&passwordReset).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *passwordResetDb) DeleteAllTokens(ctx context.Context, userId string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	err = db.Where("user_id = ?", userId).Delete(&domain.Session{}).Error
	if err != nil {
		fmt.Println("Error deleting all tokens")
		return err
	}

	return nil
}

func (p *passwordResetDb) DeleteAllTokensOlderThan15min(ctx context.Context, userId string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)
	err = db.Where("user_id = ? AND created_at < ?", userId, fifteenMinutesAgo).Delete(&domain.Session{}).Error
	if err != nil {
		fmt.Println("Error deleting all tokens older than 15 min")
		return err
	}

	return nil
}

func (p *passwordResetDb) GetAmountTokensYoungerThan15Min(ctx context.Context, userId string) (int, error) {
	db, err := getContextDb(ctx)
	if err != nil {
		return -1, err
	}

	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)

	var count int64
	err = db.Model(&domain.PasswordReset{}).Where("user_id = ? AND created_at > ?", userId, fifteenMinutesAgo).Count(&count).Error

	if err != nil {
		return -1, err
	}

	if count > 2 {
		return -1, customErrors.ErrGeneratedTooManyResetTokens
	}

	return int(count), nil
}

func (p *passwordResetDb) CheckToken(ctx context.Context, tokenString string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)

	var passwordReset domain.PasswordReset
	err = db.Model(&domain.PasswordReset{}).Where("reset_token = ? AND created_at > ? AND used = ?", tokenString, fifteenMinutesAgo, false).First(&passwordReset).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *passwordResetDb) RevokeAllTokens(ctx context.Context, userId string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	err = db.Model(&domain.PasswordReset{}).Where("user_id = ?", userId).Updates(map[string]interface{}{"used": true}).Error
	if err != nil {
		fmt.Println("Error revoking")
		return err
	}

	return nil
}
