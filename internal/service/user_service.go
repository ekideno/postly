package service

import (
	"errors"
	"github.com/ekideno/postly/internal/domain"
	"github.com/ekideno/postly/internal/security"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(user *domain.User) error {
	var err error
	user.HashedPassword, err = security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	return s.repo.Create(user)
}

func (s *UserService) GetByID(id string) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) DeleteByID(id string) error {
	return s.repo.DeleteByID(id)
}

func (s *UserService) Login(email string, password string) error {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return err
	}

	if !security.CheckPasswordHash(password, user.HashedPassword) {
		return errors.New("wrong password")
	}

	return nil
}
