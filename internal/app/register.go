package app

import (
	"bytes"
	"encoding/json"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/lib/hash"
	"github.com/ANiWarlock/gophermart/internal/models"
	"net/http"
	"strings"
)

func (a *App) RegisterHandler(rw http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var authReq models.User

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

	if err = json.Unmarshal(buf.Bytes(), &authReq); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	if authReq.Login == "" || authReq.Password == "" {
		http.Error(rw, "Wrong request format", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Login:    authReq.Login,
		Password: hash.Hash(authReq.Password),
	}

	result := a.db.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		rw.WriteHeader(http.StatusConflict)
		return
	} else if result.Error != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	newBalance := models.Balance{
		Current:   0,
		Withdrawn: 0,
		UserID:    newUser.ID,
	}

	result = a.db.Create(&newBalance)
	if result.Error != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = auth.SetCookie(newUser.ID, rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
