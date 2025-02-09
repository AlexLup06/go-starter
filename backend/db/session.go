package db

import (
	"context"

	"alexlupatsiy.com/personal-website/backend/domain"
	"alexlupatsiy.com/personal-website/backend/repository"
)

type sessionDb struct{}

func NewSessionDb() repository.SessionStorage {
	return &sessionDb{}
}

func (c *sessionDb) CreateSession(ctx context.Context, session domain.Session) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	err = db.Model(domain.Session{}).Create(&session).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *sessionDb) GetSessionById(ctx context.Context, sessionId string) (domain.Session, error) {
	db, err := getContextDb(ctx)
	if err != nil {
		return domain.Session{}, err
	}

	var session domain.Session
	err = db.Model(domain.Session{}).Where("id = ?", sessionId).First(&session).Error
	if err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

func (c *sessionDb) GetSessionByUserId(ctx context.Context, userId string) (domain.Session, error) {
	db, err := getContextDb(ctx)
	if err != nil {
		return domain.Session{}, err
	}

	var session domain.Session
	err = db.Model(domain.Session{}).Where("user_id = ?", userId).Where("revoked = ?", false).First(&session).Error
	if err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

func (c *sessionDb) DeleteSession(ctx context.Context, sessionId string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}
	err = db.Where("id = ?", sessionId).Delete(&domain.Session{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *sessionDb) RevokeSession(ctx context.Context, sessionId string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}

	err = db.Model(domain.Session{}).Where("id = ?", sessionId).Update("revoked", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *sessionDb) RevokeAllSessions(ctx context.Context, userId string) error {
	db, err := getContextDb(ctx)
	if err != nil {
		return err
	}
	err = db.Model(&domain.Session{}).Where("user_id = ?", userId).Updates(map[string]interface{}{"revoked": true}).Error
	if err != nil {
		return err
	}

	return nil
}
