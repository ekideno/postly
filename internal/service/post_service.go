package service

import (
	"errors"
	"time"

	"github.com/ekideno/postly/internal/domain"
	"github.com/ekideno/postly/internal/utils"
)

var (
	ErrInvalidUserID     = errors.New("user ID cannot be empty")
	ErrInvalidPostData   = errors.New("invalid post data")
	ErrInvalidUsername   = errors.New("username cannot be empty")
	ErrInvalidPagination = errors.New("limit and offset must be non-negative")
)

type PostService struct {
	postRepository domain.PostRepository
}

func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{
		postRepository: repo,
	}
}

func (s *PostService) Create(userID string, postReq *domain.CreatePostRequest) (*domain.Post, error) {
	if err := s.validateCreateRequest(userID, postReq); err != nil {
		return nil, err
	}

	post := s.buildPost(userID, postReq)

	if err := s.postRepository.Create(post); err != nil {
		return nil, err
	}

	if err := s.loadPostRelations(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostsByUsername(username string, limit, offset int) ([]domain.Post, error) {
	if err := s.validateUsernameAndPagination(username, limit, offset); err != nil {
		return nil, err
	}
	return s.postRepository.GetPostsByUsername(username, limit, offset)
}

func (s *PostService) GetPostsByID(userID string, limit, offset int) ([]domain.Post, error) {
	if err := s.validateUserIDAndPagination(userID, limit, offset); err != nil {
		return nil, err
	}
	return s.postRepository.GetPostsByID(userID, limit, offset)
}

func (s *PostService) GetFeed(limit, offset int) ([]domain.Post, error) {
	if err := s.validatePagination(limit, offset); err != nil {
		return nil, err
	}
	return s.postRepository.GetFeed(limit, offset)
}
func (s *PostService) validateUsernameAndPagination(username string, limit, offset int) error {
	if username == "" {
		return ErrInvalidUsername
	}
	return s.validatePagination(limit, offset)
}

func (s *PostService) validateUserIDAndPagination(userID string, limit, offset int) error {
	if userID == "" {
		return ErrInvalidUserID
	}
	return s.validatePagination(limit, offset)
}

func (s *PostService) validatePagination(limit, offset int) error {
	if limit < 0 || offset < 0 {
		return ErrInvalidPagination
	}
	return nil
}

func (s *PostService) validateCreateRequest(userID string, postReq *domain.CreatePostRequest) error {
	if userID == "" {
		return ErrInvalidUsername
	}

	if postReq == nil {
		return ErrInvalidPostData
	}

	if postReq.Title == "" && postReq.Content == "" {
		return errors.New("post must have either title or content")
	}

	return nil
}

func (s *PostService) buildPost(userID string, postReq *domain.CreatePostRequest) *domain.Post {
	post := &domain.Post{
		ID:        utils.GenerateSnowflakeID(),
		UserID:    userID,
		Title:     postReq.Title,
		Content:   postReq.Content,
		CreatedAt: time.Now(),
	}

	post.Images = s.buildPostImages(post.ID, postReq.Images)
	return post
}

func (s *PostService) buildPostImages(postID string, imageURLs []string) []domain.PostImage {
	if len(imageURLs) == 0 {
		return nil
	}

	images := make([]domain.PostImage, 0, len(imageURLs))
	for _, url := range imageURLs {
		if url != "" {
			images = append(images, domain.PostImage{
				ID:     utils.GenerateSnowflakeID(),
				PostID: postID,
				URL:    url,
			})
		}
	}

	return images
}

func (s *PostService) loadPostRelations(post *domain.Post) error {
	if err := s.postRepository.LoadAuthor(post); err != nil {
		return err
	}

	if err := s.postRepository.LoadImages(post); err != nil {
		return err
	}

	return nil
}
