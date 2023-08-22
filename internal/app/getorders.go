package app

import (
	"encoding/json"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/models"
	"net/http"
)

func (a *App) GetOrdersHandler(rw http.ResponseWriter, r *http.Request) {
	var orders []models.Order

	rw.Header().Set("Content-Type", "application/json")
	userID := auth.CurrentUser(r.Context())

	a.db.Where("user_id = ?", userID).Find(&orders).Order("created_at")
	if len(orders) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(resp)
	if err != nil {
		return
	}
}
