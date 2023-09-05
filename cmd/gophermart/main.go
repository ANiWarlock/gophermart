package main

import (
	"context"
	"errors"
	"fmt"
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
	"sync"
	"time"
)

func main() {
	const (
		timeoutServerShutdown = time.Second * 5
		timeoutShutdown       = time.Second * 10
	)
	componentsErrs := make(chan error, 1)

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	wgAccrual := &sync.WaitGroup{}
	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

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

	wg.Add(1)
	go func() {
		defer sugar.Infoln("DB closed")
		defer wg.Done()
		<-ctx.Done()
		wgAccrual.Wait()

		database.CloseDB(db)
	}()

	client, err := accrual.Init(cfg)
	if err != nil {
		sugar.Fatalf("Cannot init accrual client: %v", err)
	}

	myApp := app.NewApp(cfg, db, sugar, client)
	gmRouter := router.NewRouter(myApp, sugar)

	myApp.RestoreStatuses()
	go myApp.GetAccrual(ctx, wgAccrual)

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: gmRouter,
	}

	sugar.Infow(
		"Starting server",
		"addr", cfg.RunAddress,
	)

	go func(errs chan<- error) {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errs <- fmt.Errorf("listen and server has failed: %w", err)
		}
	}(componentsErrs)

	wg.Add(1)
	go func() {
		defer sugar.Infoln("server has been shutdown")
		defer wg.Done()
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			sugar.Errorf("an error occurred during server shutdown: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		return
	case err := <-componentsErrs:
		sugar.Errorln(err)
		cancelCtx()
	}

	go func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()

		<-ctx.Done()
		sugar.Fatal("failed to gracefully shutdown the service")
	}()
}
