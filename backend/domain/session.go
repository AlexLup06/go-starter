package domain

import "time"

type Session struct {
	DefaultFields
	UserID string `gorm:"type:uuid;column:user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	RefreshToken string    `gorm:"type:varchar(512);unique;not null;column:refresh_token"`
	IssuedAt     time.Time `gorm:"type:timestamptz;not null;default:now()"`
	ExpiresAt    time.Time `gorm:"type:timestamptz;not null"`
	Revoked      bool      `gorm:"not null;default:false"`

	UserAgent string `gorm:"type:varchar(255)"`
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) IsRevoked() bool {
	return s.Revoked
}

func (s *Session) TableName() string {
	return "session"
}
