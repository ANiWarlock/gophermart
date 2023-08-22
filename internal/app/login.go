package app

import (
	"bytes"
	"encoding/json"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"github.com/ANiWarlock/gophermart/internal/lib/hash"
	"github.com/ANiWarlock/gophermart/internal/models"
	"net/http"
)

func (a *App) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var authReq models.User
	var userFromDB models.User

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

	result := a.db.First(&userFromDB, "login = ?", authReq.Login)

	if result.Error != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if hash.Hash(authReq.Password) != userFromDB.Password {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = auth.SetCookie(userFromDB.ID, rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
