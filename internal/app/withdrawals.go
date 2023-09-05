package app

import (
	"encoding/json"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/models"
	"net/http"
)

func (a *App) GetWithdrawalsHandler(rw http.ResponseWriter, r *http.Request) {
	var withdrawals []models.Withdraw
	userID := auth.CurrentUser(r.Context())

	a.db.Where("user_id = ?", userID).Order("processed_at").Find(&withdrawals)

	if len(withdrawals) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(withdrawals)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(resp)
	if err != nil {
		return
	}
}
