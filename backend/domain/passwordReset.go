package domain

import "time"

type PasswordReset struct {
	DefaultFields
	UserID string `gorm:"type:uuid;column:user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	ExpiresAt  time.Time `gorm:"type:timestamptz;not null"`
	ResetToken string    `gorm:"type:varchar(512);unique;not null;column:reset_token"`
	Used       bool      `gorm:"not null;default:false;column:used"`
}

func (p *PasswordReset) TableName() string {
	return "password_reset"
}
