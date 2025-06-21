package database

import (
	"fmt"
	"github.com/ekideno/postly/internal/config"
	"github.com/ekideno/postly/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreDatabase struct {
	Conn *gorm.DB
}

func PostgreConnect(cfg *config.Config) (*PostgreDatabase, error) {
	dsn := getDSN(&cfg.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return &PostgreDatabase{Conn: db}, nil

}

func getDSN(db *config.Database) string {
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		db.Host, db.User, db.Password, db.Name, db.Port, db.Sslmode, db.Timezone,
	)
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{}, &domain.Post{})
}
