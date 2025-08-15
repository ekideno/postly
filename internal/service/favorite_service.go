package service

import (
	"github.com/ekideno/postly/internal/domain"
)

type FavoriteService struct {
	repo domain.FavoriteRepository
}

func NewFavoriteService(repo domain.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) AddFavorite(userID, imageID string) error {
	return s.repo.AddFavorite(userID, imageID)
}

func (s *FavoriteService) RemoveFavorite(userID, imageID string) error {
	return s.repo.RemoveFavorite(userID, imageID)
}

func (s *FavoriteService) GetFavorites(userID string, limit, offset int) ([]domain.PostImage, error) {
	return s.repo.GetFavorites(userID, limit, offset)
}
