package app

import (
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/lib/luhn"
	"github.com/ANiWarlock/gophermart/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"
	"time"
)

func (a *App) CreateOrderHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	if !luhn.Check(string(body)) {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userID := auth.CurrentUser(r.Context())
	newOrder := models.Order{
		Number:     string(body),
		Status:     models.StatusNew,
		UserID:     userID,
		UploadedAt: time.Now(),
	}

	result := a.db.Create(&newOrder)

	if err, ok := result.Error.(*pgconn.PgError); ok {
		if err.Code == pgerrcode.UniqueViolation {
			a.db.First(&newOrder, "number = ?", newOrder.Number)

			if newOrder.UserID == userID {
				rw.WriteHeader(http.StatusOK)
				return
			}
			rw.WriteHeader(http.StatusConflict)
			return
		}
	} else if result.Error != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
