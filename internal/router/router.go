package router

import (
	"github.com/ANiWarlock/gophermart/internal/app"
	"github.com/ANiWarlock/gophermart/internal/router/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(myApp *app.App, sugar *zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SugarLogger(sugar))
	r.Use(middleware.Gzip)

	r.Post("/api/user/register", myApp.RegisterHandler)
	r.Post("/api/user/login", myApp.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.CheckAuthCookie)
		r.Post("/api/user/orders", myApp.CreateOrderHandler)
		r.Get("/api/user/orders", myApp.GetOrdersHandler)
		r.Get("/api/user/balance", myApp.BalanceHandler)
		r.Post("/api/user/balance/withdraw", myApp.CreateWithdrawHandler)
		r.Get("/api/user/withdrawals", myApp.GetWithdrawalsHandler)
	})

	return r
}
