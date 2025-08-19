package domain

import (
	"time"
)

type Post struct {
	ID        string      `gorm:"primaryKey"`
	UserID    string      `gorm:"not null;index"`
	Title     string      `gorm:"not null"`
	Content   string      `gorm:"type:text"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	User      *User       `gorm:"constraint:OnDelete:CASCADE"`
	Images    []PostImage `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type PostResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	Author    PublicUserDTO `json:"author"`
	Images    []PostImage   `json:"images"`
}

type PostRepository interface {
	Create(post *Post) error
	GetPostsByUsername(username string, limit, offset int) ([]Post, error)
	GetPostsByID(userID string, limit, offset int) ([]Post, error)
	LoadAuthor(post *Post) error
	GetFeed(limit, offset int) ([]Post, error)
	LoadImages(post *Post) error
}

type CreatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type PostImage struct {
	ID     string `gorm:"primaryKey" json:"id"`
	PostID string `gorm:"not null;index" json:"post_id"`
	URL    string `gorm:"not null" json:"url"`
}
