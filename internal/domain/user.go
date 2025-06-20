package domain

import "gorm.io/gorm"

type User struct {
	ID             string `gorm:"primaryKey"`
	Username       string `gorm:"uniqueIndex;not null" json:"username"`
	Email          string `gorm:"uniqueIndex;not null" json:"email"`
	Password       string `gorm:"-" json:"-"`
	HashedPassword string `gorm:"not null" json:"-"`
	gorm.Model     `json:"-"`
}
type PublicUserDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	DeleteByID(id string) error
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}
