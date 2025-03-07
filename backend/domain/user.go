package domain

type User struct {
	DefaultFields
	Username       string          `gorm:"column:username"`
	Email          *string         `gorm:"column:email"`
	AuthProviders  []AuthProvider  `gorm:"foreignKey:UserID"`
	Sessions       []Session       `gorm:"foreignKey:UserID"`
	PasswordResets []PasswordReset `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "user"
}
