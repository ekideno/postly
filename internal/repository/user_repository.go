package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/ekideno/postly/internal/domain"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepository struct {
	db *gorm.DB
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}
}

func getDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)
}

func NewUserRepository() (*UserRepository, error) {
	loadEnv()

	dsn := getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
	log.Printf("Deleting user with id: %s\n", id)
	result := r.db.Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}
	return nil
}
