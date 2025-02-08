package domain

import "time"

type Session struct {
	DefaultFields
	ExpiresAt time.Time `gorm:"column:expires_at"`
	// Either UserId or TeamId is set. Not both.
	UserId *string `gorm:"type:uuid;column:user_id;default:null"`
	User   *User   `gorm:"foreignKey:UserId"`
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) TableName() string {
	return "session"
}
