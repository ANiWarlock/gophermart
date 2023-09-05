package middleware

import (
	"context"
	"github.com/ANiWarlock/gophermart/internal/lib/auth"
	"net/http"
)

func CheckAuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("auth")
		if err != nil {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := auth.GetUserID(tokenString.Value)
		if err != nil {
			http.Error(rw, "500", http.StatusInternalServerError)
			return
		}
		if userID == 0 {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = setCtxUserID(r, userID)
		next.ServeHTTP(rw, r)
	})
}

func setCtxUserID(r *http.Request, userID uint) *http.Request {
	ctx := context.WithValue(r.Context(), auth.CtxKeyUserID, userID)
	return r.WithContext(ctx)
}
