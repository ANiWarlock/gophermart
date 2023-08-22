package app

import (
	"github.com/ANiWarlock/gophermart/cmd/gophermart/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg   *config.AppConfig
	db    *gorm.DB
	sugar *zap.SugaredLogger
}

func NewApp(cfg *config.AppConfig, db *gorm.DB, sugar *zap.SugaredLogger) *App {
	return &App{
		cfg:   cfg,
		db:    db,
		sugar: sugar,
	}
}
