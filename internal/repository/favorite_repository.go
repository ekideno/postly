package repository

import (
	"fmt"

	"github.com/ekideno/postly/internal/domain"
	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) (*FavoriteRepository, error) {
	return &FavoriteRepository{db: db}, nil
}

func (r *FavoriteRepository) AddFavorite(userID, imageID string) error {
	var user domain.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var image domain.PostImage
	if err := r.db.First(&image, "id = ?", imageID).Error; err != nil {
		return fmt.Errorf("image not found")
	}

	return r.db.Model(&user).Association("Favorites").Append(&image)
}

func (r *FavoriteRepository) RemoveFavorite(userID, imageID string) error {
	var user domain.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var image domain.PostImage
	if err := r.db.First(&image, "id = ?", imageID).Error; err != nil {
		return fmt.Errorf("image not found")
	}

	return r.db.Model(&user).Association("Favorites").Delete(&image)

}

func (r *FavoriteRepository) GetFavorites(userID string, limit, offset int) ([]domain.PostImage, error) {
	var images []domain.PostImage
	err := r.db.
		Joins("JOIN favorites f ON f.image_id = post_images.id").
		Where("f.user_id = ?", userID).
		Order("post_images.id").
		Limit(limit).
		Offset(offset).
		Find(&images).Error
	return images, err
}
