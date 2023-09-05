package app

import (
	"context"
	"errors"
	"github.com/ANiWarlock/gophermart/internal/lib/accrual"
	"github.com/ANiWarlock/gophermart/internal/models"
	"gorm.io/gorm"
	"math"
	"sync"
	"time"
)

func (a *App) RestoreStatuses() {
	// вернем статусы в NEW, если сервис упал в процессе
	a.db.Where("status IN ?", []string{models.StatusProcessing}).Update("status", models.StatusNew)
}

func (a *App) GetAccrual(ctx context.Context, wg *sync.WaitGroup) {
	var orders []models.Order

	for {
		select {
		case <-ctx.Done():
			return
		default:
			result := a.db.Where("status IN ?", []string{models.StatusNew}).Find(&orders)

			if result.RowsAffected == 0 {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, order := range orders {
				order.Status = models.StatusProcessing
				a.db.Save(&order)
				wg.Add(1)
				go a.process(ctx, wg, order)
			}
		}
	}
}

func (a *App) process(ctx context.Context, wg *sync.WaitGroup, order models.Order) {
	var balance models.Balance
	var tries = 0

get:
	for {
		select {
		case <-ctx.Done():
			order.Status = models.StatusNew
			a.db.Save(&order)
			wg.Done()
			return
		default:
			if tries > 10 {
				order.Status = models.StatusNew
				a.db.Save(&order)
				wg.Done()
				return
			}
			accr, err := a.client.Get(order.Number)

			if err != nil && errors.Is(err, accrual.ErrTooManyRequests) {
				a.sugar.Error("Accrual.GET too many requests, sleep 5 seconds")
				time.Sleep(5 * time.Second)
				continue
			} else if err != nil {
				a.sugar.Errorf("Accrual.GET error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if accr == nil {
				order.Status = models.StatusInvalid
				a.db.Save(&order)
				wg.Done()
				break get
			}

			switch accr.Status {
			case models.StatusRegistered:
				tries++
				time.Sleep(1 * time.Second)
			case models.StatusProcessing:
				tries++
				time.Sleep(1 * time.Second)
			case models.StatusInvalid:
				order.Status = models.StatusInvalid
				a.db.Save(&order)
			case models.StatusProcessed:
				err = a.db.Transaction(func(tx *gorm.DB) error {
					order.Status = models.StatusProcessed
					order.Accrual = accr.Accrual
					result := tx.Save(&order)
					if result.Error != nil {
						a.sugar.Errorf("failed to save order: %v", result.Error)
						return result.Error
					}

					a.db.First(&balance, "user_id = ?", order.UserID)

					balance.Current = math.Round((balance.Current+accr.Accrual)*100) / 100

					result = tx.Save(&balance)
					if result.Error != nil {
						a.sugar.Errorf("failed to save balance: %v", result.Error)
						return result.Error
					}
					return nil
				})
				wg.Done()
				if err != nil {
					a.sugar.Errorf("failed to process accrual: %v", err)
				}
				return
			}
		}
	}
}
