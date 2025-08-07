package domain

type User struct {
	ID             string  `gorm:"primaryKey"`
	Username       string  `gorm:"uniqueIndex;not null"`
	Email          string  `gorm:"uniqueIndex;not null"`
	Password       string  `gorm:"-"`
	HashedPassword string  `gorm:"not null"`
	Posts          []Post  `gorm:"foreignKey:UserID"`
	Bio            string  `json:"bio"`
	AvatarURL      string  `json:"avatar_url"`
	BannerURL      string  `json:"banner_url"`
	Following      []*User `gorm:"many2many:user_followings;joinForeignKey:UserID;joinReferences:FollowingID"`
	Followers      []*User `gorm:"many2many:user_GetByUsernamefollowings;joinForeignKey:FollowingID;joinReferences:UserID"`
}

type PublicUserDTO struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	BannerURL   string `json:"banner_url"`
	IsFollowing bool   `json:"is_following"`
	IsMe        bool   `json:"is_me"`
}

type PrivateUserDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	BannerURL string `json:"banner_url"`
}

type UpdateUserDTO struct {
	Email     *string `json:"email,omitempty"`
	Username  *string `json:"username,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL string  `json:"avatar_url"`
	BannerURL string  `json:"banner_url"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	DeleteByID(id string) error
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Follow(userID string, targetID string) error
	GetFollowing(userID string) ([]User, error)
	IsFollowing(followerID, followingID string) (bool, error)
	Unfollow(userID string, targetID string) error
	GetFollowedUserIDs(userID string, authorIDs []string) ([]string, error)
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
