package main

import (
	"log"
	"net/http"
	"userwalletservice/internal/infrastructure/database"
	"userwalletservice/internal/infrastructure/jwt"
	"userwalletservice/internal/router"
	"userwalletservice/internal/router/middleware"

	// Слой репозиториев
	userRepo "userwalletservice/internal/repository/user"
	walletRepo "userwalletservice/internal/repository/wallet"

	// Слой сервисов и контроллеров
	userCtrl "userwalletservice/internal/controller/user"
	walletCtrl "userwalletservice/internal/controller/wallet"
	userServ "userwalletservice/internal/service/user"
)

func main() {
	// 1. Подключение к БД
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// 2. Инициализируем JWT Менеджер (пока хардкодим ключ, потом вынесешь в конфиг)
	jwtSecret := "super-secret-wallet-key-2026"
	jwtManager := jwt.New(jwtSecret)

	// 3. Собираем слой репозиториев и сервисов
	userRepository := userRepo.New(db)
	walletRepository := walletRepo.New(db)

	// 🔥 3. Передаем walletRepository третьим аргументом в конструктор сервиса
	userService := userServ.New(userRepository, jwtManager, walletRepository)

	// 4. Собираем слой контроллеров и мидлваров
	userController := userCtrl.New(userService)
	walletController := walletCtrl.New(userService)
	authMiddleware := middleware.New(jwtManager)

	// 5. Инициализируем роутер gorilla/mux и настраиваем маршруты
	r := router.New(userController, walletController, authMiddleware)
	muxRouter := r.InitRoutes() // Это вернет нам настроенный *mux.Router

	log.Println("🚀 Server started on :8080")

	// 6. Запуск сервера
	if err := http.ListenAndServe(":8080", muxRouter); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
