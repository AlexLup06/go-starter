package domain

type AuthProvider struct {
	DefaultFields
	UserID         uint    `gorm:"index"`
	Provider       string  `gorm:"type:varchar(50);not null"`
	ProviderUserID *string `gorm:"uniqueIndex:provider_user_idx"`
	PasswordHash   *string `gorm:"type:varchar(255)"`
}
