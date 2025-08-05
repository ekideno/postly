package service

import (
	"time"

	"github.com/ekideno/postly/internal/domain"
	"github.com/ekideno/postly/internal/utils"
)

type PostService struct {
	postRepository domain.PostRepository
}

func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{postRepository: repo}
}

func (s *PostService) Create(userID string, postReq *domain.CreatePostRequest) (*domain.Post, error) {
	post := &domain.Post{
		ID:        utils.GenerateSnowflakeID(),
		UserID:    userID,
		Title:     postReq.Title,
		Content:   postReq.Content,
		CreatedAt: time.Now(),
	}

	var images []domain.PostImage
	for _, url := range postReq.Images {
		images = append(images, domain.PostImage{
			ID:     utils.GenerateSnowflakeID(),
			PostID: post.ID,
			URL:    url,
		})
	}

	post.Images = images

	err := s.postRepository.Create(post)
	if err != nil {
		return nil, err
	}

	err = s.postRepository.LoadAuthor(post)
	if err != nil {
		return nil, err
	}

	err = s.postRepository.LoadImages(post)
	return post, nil
}

func (s *PostService) GetPostsByUsername(username string, limit, offset int) ([]domain.Post, error) {
	return s.postRepository.GetPostsByUsername(username, limit, offset)
}

func (s *PostService) GetPostsByID(userID string, limit, offset int) ([]domain.Post, error) {
	return s.postRepository.GetPostsByID(userID, limit, offset)
}

func (s *PostService) GetFeed(limit, offset int) ([]domain.Post, error) {
	return s.postRepository.GetFeed(limit, offset)
}
func (s *PostService) LoadImages(post *domain.Post) error {
	return s.LoadImages(post)
}
