package app

import (
	"encoding/json"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/models"
	"net/http"
)

func (a *App) BalanceHandler(rw http.ResponseWriter, r *http.Request) {
	var balance models.Balance

	userID := auth.CurrentUser(r.Context())

	a.db.First(&balance, "user_id = ?", userID)

	resp, err := json.Marshal(balance)
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
