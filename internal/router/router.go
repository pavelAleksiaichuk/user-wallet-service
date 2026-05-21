package router

import (
	userCntrl "userwalletservice/internal/controller/user"
	walletCntrl "userwalletservice/internal/controller/wallet"
	"userwalletservice/internal/router/middleware"

	"github.com/gorilla/mux"
)

type Router struct {
	// Контроллеры из полей структуры убираем — теперь они передаются локально в методы
	authMid *middleware.AuthMiddleware
}

// 1. Конструктор стал чистым. Принимает ТОЛЬКО мидлвар
func New(authMid *middleware.AuthMiddleware) *Router {
	return &Router{
		authMid: authMid,
	}
}

// 2. InitRoutes теперь принимает контроллеры на лету
func (r *Router) InitRoutes(userController *userCntrl.Controller, walletController *walletCntrl.Controller) *mux.Router {
	// Как и в твоем старом коде — создаем базовый роутер прямо здесь
	mainRouter := mux.NewRouter()

	// Создаем общую группу для API v1
	v1 := mainRouter.PathPrefix("/api/v1").Subrouter()

	// Создаем один суб-роутер для ВСЕХ защищенных эндпоинтов v1
	protected := v1.PathPrefix("/").Subrouter()
	protected.Use(r.authMid.Handler)

	// 🔥 Вызываем функции, как просил Костя:
	r.NewUser(v1, protected, userController)
	r.NewWallet(protected, walletController)

	return mainRouter
}

// 3. Функция для сущности User
func (r *Router) NewUser(public *mux.Router, protected *mux.Router, c *userCntrl.Controller) {
	public.HandleFunc("/auth/register", c.RegisterUser).Methods("POST")
	public.HandleFunc("/auth/login", c.LoginUser).Methods("POST")
	protected.HandleFunc("/auth/profile", c.GetProfile).Methods("GET")
}

// 4. Функция для сущности Wallet
func (r *Router) NewWallet(protected *mux.Router, c *walletCntrl.Controller) {
	protected.HandleFunc("/account/balance", c.GetBalance).Methods("GET")
	protected.HandleFunc("/account/deposit", c.Deposit).Methods("POST")
	protected.HandleFunc("/account/withdraw", c.Withdraw).Methods("POST")
	protected.HandleFunc("/account/transfer", c.Transfer).Methods("POST")
}
