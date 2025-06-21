package domain

type User struct {
	ID             string `gorm:"primaryKey"`
	Username       string `gorm:"uniqueIndex;not null"`
	Email          string `gorm:"uniqueIndex;not null"`
	Password       string `gorm:"-"`
	HashedPassword string `gorm:"not null"`

	Posts []Post `gorm:"foreignKey:UserID"`
}

type PublicUserDTO struct {
	ID       string `json:"id"`
	Username string `json:"name"`
}

type PrivateUserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	DeleteByID(id string) error
	GetByUsername(username string) (*User, error)
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
