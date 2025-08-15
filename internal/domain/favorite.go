package domain

import "time"

type Favorite struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index;not null"`
	PostID    string `gorm:"index;not null"`
	CreatedAt time.Time
}
type FavoriteRepository interface {
	AddFavorite(userID, imageID string) error
	RemoveFavorite(userID, imageID string) error
	GetFavorites(userID string, limit, offset int) ([]PostImage, error)
}
