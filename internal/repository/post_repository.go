package repository

import (
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

func (r *PostRepository) GetByUserID(userID string, limit, offset int) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.
		Where("user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, err
}
