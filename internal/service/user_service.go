package service

import (
	"errors"
	"github.com/ekideno/postly/internal/domain"
	"github.com/ekideno/postly/internal/security"
	"github.com/ekideno/postly/internal/utils"
)

type UserService struct {
	repo       domain.UserRepository
	jwtManager *security.JWTManager
}

func NewUserService(repo domain.UserRepository, jwtManager *security.JWTManager) *UserService {
	return &UserService{repo: repo, jwtManager: jwtManager}
}

func (s *UserService) Register(user *domain.User) (string, error) {
	user.ID = utils.GenerateSnowflakeID()

	var err error
	user.HashedPassword, err = security.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	err = s.repo.Create(user)
	if err != nil {
		return "", err
	}
	return s.jwtManager.GenerateToken(user.ID, user.Username)
}

func (s *UserService) GetByID(id string) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) DeleteByID(id string) error {
	return s.repo.DeleteByID(id)
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if !security.CheckPasswordHash(password, user.HashedPassword) {
		return "", errors.New("wrong password")
	}

	return s.jwtManager.GenerateToken(user.ID, user.Username)
}

func (s *UserService) GetByUsername(username string) (*domain.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) UpdateUserProfile(userID string, dto domain.UpdateUserDTO) (*domain.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.Bio != nil {
		user.Bio = *dto.Bio
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateAvatar(userID string, path string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	user.AvatarURL = path
	return s.repo.Update(user)
}

func (s *UserService) UpdateBanner(userID string, path string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	user.BannerURL = path
	return s.repo.Update(user)
}
