package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey            string `env:"SECRET_KEY"`
}

// адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
//адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
//адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.

func InitConfig() (*AppConfig, error) {
	cfg := AppConfig{}
	cfg.parseFlags()
	err := cfg.parseEnvs()
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

func (c *AppConfig) parseFlags() {
	flag.StringVar(&c.RunAddress, "a", "localhost:8080", "адрес и порт запуска сервиса")
	flag.StringVar(&c.DatabaseURI, "d", "host=localhost user=gopher password=123 dbname=gophermart sslmode=disable", "адрес подключения к базе данных")
	flag.StringVar(&c.AccrualSystemAddress, "r", "http://localhost:8081", "адрес системы расчёта начислений")
	flag.StringVar(&c.SecretKey, "s", "", "Ключ шифрования")
	flag.Parse()
}

func (c *AppConfig) parseEnvs() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("failed to parse env vars: %w", err)
	}
	return nil
}
