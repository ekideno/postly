package repository

import (
	"fmt"

	"github.com/ekideno/postly/internal/domain"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) (*PostRepository, error) {
	return &PostRepository{db: db}, nil
}

func (r *PostRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) LoadAuthor(post *domain.Post) error {
	return r.db.Preload("User").First(post, "id = ?", post.ID).Error
}

func (r *PostRepository) GetPostsByUsername(username string, limit, offset int) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.
		Joins("JOIN users ON users.id = posts.user_id").
		Where("users.username = ?", username).
		Preload("User").
		Preload("Images").
		Order("posts.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, err
}

func (r *PostRepository) GetPostsByID(userID string, limit, offset int) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.
		Joins("JOIN users ON users.id = posts.user_id").
		Where("users.id = ?", userID).
		Preload("User").
		Preload("Images").
		Order("posts.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, err
}

func (r *PostRepository) GetFeed(limit, offset int) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.
		Preload("User").
		Preload("Images").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, err
}

func (r *PostRepository) LoadImages(post *domain.Post) error {
	r.db.Preload("User").Preload("Images").First(post, "id = ?", post.ID)
	fmt.Println(post)
	return nil

}
