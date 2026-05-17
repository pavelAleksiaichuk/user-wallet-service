package router

import (
	userCntrl "userwalletservice/internal/controller/user"
	"userwalletservice/internal/router/middleware"

	"github.com/gorilla/mux"
)

type Router struct {
	userController *userCntrl.UserController
	authMid  *middleware.AuthMiddleware
}

func New(userController *userCntrl.UserController, authMid *middleware.AuthMiddleware) *Router {
	return &Router{
		userController: userController,
		authMid:  authMid,
	}
}

// InitRoutes настраивает gorilla/mux роутер
func (r *Router) InitRoutes() *mux.Router {
	// Главный роутер
	mainRouter := mux.NewRouter()

	// Создаем общую группу для API v1. 
	// Теперь не нужно писать "/api/v1" в каждом HandleFunc!
	v1 := mainRouter.PathPrefix("/api/v1").Subrouter()

	// ==========================================
	// 1. ПУБЛИЧНЫЕ РОУТЫ (Группа /api/v1/auth)
	// ==========================================
	authPublic := v1.PathPrefix("/auth").Subrouter()
	authPublic.HandleFunc("/register", r.userController.RegisterUser).Methods("POST")
	authPublic.HandleFunc("/login", r.userController.LoginUser).Methods("POST")

	// ==========================================
	// 2. ЗАЩИЩЕННЫЕ РОУТЫ (Группа /api/v1/auth)
	// ==========================================
	// Создаем отдельный суб-роутер для защищенных методов auth
	authProtected := v1.PathPrefix("/auth").Subrouter()
	
	// Вешаем мидлвар на ВСЕ роуты внутри этого суб-роутера автоматически!
	authProtected.Use(r.authMid.Handler) 
	
	// URL будет: /api/v1/auth/profile. И никаких ручных оберток!
	authProtected.HandleFunc("/profile", r.userController.GetProfile).Methods("GET")

	// ==========================================
	// ПРИМЕР НА БУДУЩЕЕ: Роуты кошелька
	// ==========================================
	// walletProtected := v1.PathPrefix("/wallets").Subrouter()
	// walletProtected.Use(r.authMid.Handler)
	// walletProtected.HandleFunc("", r.walletCtrl.CreateWallet).Methods("POST") // /api/v1/wallets
	// walletProtected.HandleFunc("/{id}/balance", r.walletCtrl.GetBalance).Methods("GET") // /api/v1/wallets/1/balance

	return mainRouter
}
