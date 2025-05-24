package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	DeleteByID(id string) error
}
