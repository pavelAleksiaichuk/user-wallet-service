package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"userwalletservice/internal/infrastructure/database"
	"userwalletservice/internal/infrastructure/jwt"
	"userwalletservice/internal/router"
	"userwalletservice/internal/router/middleware"

	// Слой репозиториев
	userRepository "userwalletservice/internal/repository/user"
	walletRepository "userwalletservice/internal/repository/wallet"

	// Слой контроллеров
	userController "userwalletservice/internal/controller/user"
	walletController "userwalletservice/internal/controller/wallet"

	// Слой сервисов
	userService "userwalletservice/internal/service/user"
	walletService "userwalletservice/internal/service/wallet"
)

func main() {
	// 1. Подключение к БД
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 2. Инициализируем JWT Менеджер (пока хардкодим ключ, потом вынесешь в конфиг)
	jwtSecret := "super-secret-wallet-key-2026"
	jwtManager := jwt.New(jwtSecret)

	// 3. Собираем слой репозиториев и сервисов
	userRepository := userRepository.New(db)
	walletRepository := walletRepository.New(db)

	// 🔥 3. Передаем walletRepository третьим аргументом в конструктор сервиса
	userService := userService.New(userRepository, jwtManager, walletRepository)
	walletService := walletService.New(walletRepository)

	// 4. Собираем слой контроллеров и мидлваров
	userController := userController.New(userService)
	walletController := walletController.New(walletService)
	authMiddleware := middleware.New(jwtManager)

	// 5. Инициализируем роутер gorilla/mux и настраиваем маршруты
	r := router.New(authMiddleware)
	muxRouter := r.InitRoutes(userController, walletController) // Это вернет нам настроенный *mux.Router

	// 6. Настраиваем конфигурацию HTTP-сервера
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      muxRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Создаем контекст, который отменится при системных сигналах (Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("🚀 Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	// Ожидаем системный сигнал завершения (код замрёт на этой строке)
	<-ctx.Done()
	log.Println("🛑 Shutting down gracefully...")

	// 9. Даем серверу 30 секунд на то, чтобы завершить активные сетевые запросы
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// 10. Только ПОСЛЕ остановки сервера безопасно закрываем соединение с БД
	if err := db.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	log.Println("✅ Server stopped successfully")
}
