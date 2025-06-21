package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/ekideno/postly/internal/config"
	"github.com/ekideno/postly/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepository struct {
	db *gorm.DB
}

func getDSN(db *config.Database) string {
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		db.Host, db.User, db.Password, db.Name, db.Port, db.Sslmode, db.Timezone,
	)
}

func NewUserRepository(cfg *config.Config) (*UserRepository, error) {
	dsn := getDSN(&cfg.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}

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
