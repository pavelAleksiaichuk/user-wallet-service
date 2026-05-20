package router

import (
	userCntrl "userwalletservice/internal/controller/user"
	walletCntrl "userwalletservice/internal/controller/wallet" // 🔥 1. Импортируем контроллер кошелька
	"userwalletservice/internal/router/middleware"

	"github.com/gorilla/mux"
)

type Router struct {
	userController   *userCntrl.UserController
	walletController *walletCntrl.WalletController // 🔥 2. Добавляем поле для кошелька
	authMid          *middleware.AuthMiddleware
}

// 🔥 3. Обновляем конструктор: принимает userController, walletController и authMid
func New(
	userController *userCntrl.UserController,
	walletController *walletCntrl.WalletController,
	authMid *middleware.AuthMiddleware,
) *Router {
	return &Router{
		userController:   userController,
		walletController: walletController,
		authMid:          authMid,
	}
}

// InitRoutes настраивает gorilla/mux роутер
func (r *Router) InitRoutes() *mux.Router {
	mainRouter := mux.NewRouter()

	// Создаем общую группу для API v1
	v1 := mainRouter.PathPrefix("/api/v1").Subrouter()

	// ==========================================
	// 1. ПУБЛИЧНЫЕ РОУТЫ (Без мидлвара)
	// ==========================================
	v1.HandleFunc("/auth/register", r.userController.RegisterUser).Methods("POST")
	v1.HandleFunc("/auth/login", r.userController.LoginUser).Methods("POST")

	// ==========================================
	// 2. ЗАЩИЩЕННЫЕ РОУТЫ (Через отдельный Subrouter с мидлваром)
	// ==========================================
	// Создаем один суб-роутер для ВСЕХ защищенных эндпоинтов v1
	protected := v1.PathPrefix("/").Subrouter()
	protected.Use(r.authMid.Handler) // Защищаем ВСЁ, что добавим ниже

	// Теперь просто перечисляем защищенные эндпоинты:
	protected.HandleFunc("/auth/profile", r.userController.GetProfile).Methods("GET")
	protected.HandleFunc("/account/balance", r.walletController.GetBalance).Methods("GET")
	protected.HandleFunc("/account/deposit", r.walletController.Deposit).Methods("POST")
	protected.HandleFunc("/account/withdraw", r.walletController.Withdraw).Methods("POST")
	protected.HandleFunc("/account/transfer", r.walletController.Transfer).Methods("POST")

	return mainRouter
}
