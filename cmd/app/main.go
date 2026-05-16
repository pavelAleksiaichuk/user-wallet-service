package main

import (
	"log"
	"net/http"

	"userwalletservice/internal/controller"
	"userwalletservice/internal/infrastructure" // Для JWTManager
	"userwalletservice/internal/infrastructure/database"
	"userwalletservice/internal/repository"
	"userwalletservice/internal/router" // Для роутера и мидлвара
	"userwalletservice/internal/service"
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
	jwtManager := infrastructure.NewJWTManager(jwtSecret)

	// 3. Собираем слой репозиториев и сервисов (передаем jwtManager в UserService)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)

	// 4. Собираем слой контроллеров и мидлваров
	userController := controller.NewUserController(userService)
	authMiddleware := router.NewAuthMiddleware(jwtManager)

	// 5. Инициализируем роутер gorilla/mux и настраиваем маршруты
	r := router.NewRouter(userController, authMiddleware)
	muxRouter := r.InitRoutes() // Это вернет нам настроенный *mux.Router

	log.Println("🚀 Server started on :8080")

	// 6. Вместо nil передаем наш muxRouter, чтобы gorilla/mux управлял запросами
	if err := http.ListenAndServe(":8080", muxRouter); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}