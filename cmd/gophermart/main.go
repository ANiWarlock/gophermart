package main

import (
	"context"
	"github.com/ANiWarlock/gophermart/cmd/gophermart/config"
	"github.com/ANiWarlock/gophermart/cmd/logger"
	"github.com/ANiWarlock/gophermart/internal/app"
	"github.com/ANiWarlock/gophermart/internal/database"
	"github.com/ANiWarlock/gophermart/internal/lib/accrual"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sugar, err := logger.Initialize("info")
	if err != nil {
		log.Fatalf("Cannot init logger: %v", err)
	}

	cfg, err := config.InitConfig()
	if err != nil {
		sugar.Fatalf("Cannot init config: %v", err)
	}
	auth.SetSecretKey(cfg)

	db, err := database.InitDB(*cfg)
	if err != nil {
		sugar.Fatalf("Cannot init database: %v", err)
	}
	defer database.CloseDB(db)

	myApp := app.NewApp(cfg, db, sugar)
	gmRouter := router.NewRouter(myApp, sugar)

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: gmRouter,
	}

	sugar.Infow(
		"Starting server",
		"addr", cfg.RunAddress,
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("listen: %s\n", err)
		}
	}()

	accrual.Init(cfg)
	go myApp.GetAccrual()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Graceful shutdown: start (5 sec)")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Graceful shutdown: error: %s", err)
	}

	if <-ctx.Done(); true {
		sugar.Info("Graceful shutdown: timed out.")
	}
}
