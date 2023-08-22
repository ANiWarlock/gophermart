package database

import (
	"fmt"
	"github.com/ANiWarlock/gophermart/cmd/gophermart/config"
	"github.com/ANiWarlock/gophermart/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

func InitDB(conf config.AppConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(conf.DatabaseURI), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to init database database: %w", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	err = db.AutoMigrate(&models.Withdraw{})
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	err = db.AutoMigrate(&models.Balance{})
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}

func CloseDB(db *gorm.DB) {
	dbInstance, _ := db.DB()
	_ = dbInstance.Close()
}
