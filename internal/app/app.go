package app

import (
	"github.com/ANiWarlock/gophermart/cmd/gophermart/config"
	"github.com/ANiWarlock/gophermart/internal/lib/accrual"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.AppConfig
	db     *gorm.DB
	sugar  *zap.SugaredLogger
	client *accrual.Client
}

func NewApp(cfg *config.AppConfig, db *gorm.DB, sugar *zap.SugaredLogger, client *accrual.Client) *App {
	return &App{
		cfg:    cfg,
		db:     db,
		sugar:  sugar,
		client: client,
	}
}
