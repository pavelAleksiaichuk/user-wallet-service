package main

import (
	"fmt"
	"log"
	"net/http"

	"userwalletservice/internal/controller"
	"userwalletservice/internal/infrastructure/database"
	"userwalletservice/internal/repository"
	"userwalletservice/internal/service"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	http.HandleFunc("/auth/register", userController.RegisterUser)
	http.HandleFunc("/auth/login", userController.LoginUser)
	http.HandleFunc("/user", userController.GetUserByID)

	fmt.Println("🚀 Server started on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
