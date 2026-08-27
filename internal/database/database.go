package database

import (
	"errors"

	"github.com/angyxys/kexel/internal/config"
	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.Config) (*gorm.DB, error) {
	dsn := cfg.DATABASE_DSN
	if dsn == "" {
		return nil, errors.New("missing database dsn")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Player{})
	return db, nil
}
