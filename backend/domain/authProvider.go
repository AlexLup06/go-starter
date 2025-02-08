package domain

type AuthProvider struct {
	DefaultFields
	UserID         string  `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Method         string  `gorm:"type:varchar(50);not null"`
	ProviderUserID *string `gorm:"type:varchar(255);uniqueIndex:unique_provider_user,priority:2"`
	PasswordHash   *string `gorm:"type:varchar(255)"`
}

func (AuthProvider) TableName() string {
	return "auth_provider"
}
