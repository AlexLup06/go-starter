package domain

type User struct {
	DefaultFields
	Username      string         `gorm:"column:username"`
	email         *string        `gorm:"column:email"`
	AuthProviders []AuthProvider `gorm:"foreignKey:UserID"`
}
