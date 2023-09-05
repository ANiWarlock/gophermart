package app

import (
	"bytes"
	"encoding/json"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/lib/luhn"
	"github.com/ANiWarlock/gophermart/internal/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (a *App) CreateWithdrawHandler(rw http.ResponseWriter, r *http.Request) {
	var (
		buf      bytes.Buffer
		withdraw models.Withdraw
		balance  models.Balance
	)
	userID := auth.CurrentUser(r.Context())

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	if buf.String() == "" {
		http.Error(rw, "Empty body!", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &withdraw); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	if withdraw.Order == "" {
		http.Error(rw, "Wrong request format", http.StatusUnprocessableEntity)
		return
	}

	if !luhn.Check(withdraw.Order) {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = a.db.Transaction(func(tx *gorm.DB) error {
		tx.First(&balance, "user_id = ?", userID)

		if balance.Current >= withdraw.Sum {
			balance.Current -= withdraw.Sum
			balance.Withdrawn += withdraw.Sum
			result := tx.Save(&balance)

			if result.Error != nil {
				a.sugar.Errorf("failed to save user: %v", result.Error)
				return result.Error
			}

			withdraw.UserID = userID
			withdraw.ProcessedAt = time.Now()
			result = tx.Save(&withdraw)

			if result.Error != nil {
				a.sugar.Errorf("failed to save withdraw: %v", result.Error)
				return result.Error
			}

			rw.WriteHeader(http.StatusOK)
			return nil
		} else {
			rw.WriteHeader(http.StatusPaymentRequired)
			return nil
		}
	})

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		a.sugar.Errorf("failed to withdraw: %v", err)
	}
}
