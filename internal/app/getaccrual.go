package app

import (
	"github.com/ANiWarlock/gophermart/internal/lib/accrual"
	"github.com/ANiWarlock/gophermart/internal/models"
	"gorm.io/gorm"
	"math"
	"time"
)

func (a *App) GetAccrual() {
	var orders []models.Order
	for {
		result := a.db.Where("status IN ?", []string{"NEW"}).Find(&orders)

		if result.RowsAffected == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		for _, order := range orders {
			a.process(&order)
		}
	}
}

func (a *App) process(order *models.Order) {
	var balance models.Balance
	var tries = 0

	order.Status = "PROCESSING"
	a.db.Save(&order)

get:
	for {
		if tries > 5 {
			break
		}
		accr, err := accrual.Get(order.Number)

		if err != nil {
			a.sugar.Errorf("Accrual.GET error: %v", err)
			time.Sleep(1 * time.Second)
			tries++
			continue
		}

		if accr == nil {
			order.Status = "INVALID"
			a.db.Save(&order)
			break
		}

		switch accr.Status {
		case "REGISTERED":
			time.Sleep(1 * time.Second)
			continue
		case "PROCESSING":
			time.Sleep(1 * time.Second)
			continue
		case "INVALID":
			order.Status = "INVALID"
			a.db.Save(&order)
			break get
		case "PROCESSED":
			err = a.db.Transaction(func(tx *gorm.DB) error {
				order.Status = "PROCESSED"
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
			if err != nil {
				a.sugar.Errorf("failed to process accrual: %v", err)
			}
			break get
		}
	}
}
