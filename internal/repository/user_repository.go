package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/ekideno/postly/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) (*UserRepository, error) {
	return &UserRepository{
		db: db,
	}, nil
}

func (r UserRepository) Create(user *domain.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r UserRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) DeleteByID(id string) error {
	log.Printf("Deleting user with id: %v\n", id)
	result := r.db.Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %v", id)
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no user found with email %v", email)
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no user found with username %v", username)
		}
	}
	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Follow(userID string, targetID string) error {
	var user, target domain.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}
	if err := r.db.First(&target, "id = ?", targetID).Error; err != nil {
		return err
	}
	return r.db.Model(&user).Association("Following").Append(&target)
}

func (r *UserRepository) GetFollowing(userID string) ([]domain.User, error) {
	var user domain.User

	err := r.db.Preload("Following").First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	following := make([]domain.User, len(user.Following))
	for i, u := range user.Following {
		following[i] = *u
	}

	return following, nil
}

func (r *UserRepository) Unfollow(userID string, targetID string) error {
	var user, target domain.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	if err := r.db.First(&target, "id = ?", targetID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Following").Delete(&target)
}

func (r *UserRepository) IsFollowing(followerID, followingID string) (bool, error) {
	var count int64
	err := r.db.Table("user_followings").
		Where("user_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	fmt.Println(count > 0)
	return count > 0, nil
}

func (r *UserRepository) GetFollowedUserIDs(userID string, authorIDs []string) ([]string, error) {
	var followedIDs []string

	err := r.db.
		Table("user_followings").
		Where("user_id = ? AND following_id IN ?", userID, authorIDs).
		Pluck("following_id", &followedIDs).Error

	if err != nil {
		return nil, err
	}
	return followedIDs, nil
}
