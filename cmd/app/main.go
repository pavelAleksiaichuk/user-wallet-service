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

	// DB
	database.Connect()

	// DI
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			userController.CreateUser(w, r)

		case http.MethodGet:
			userController.GetUser(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("🚀 Server started on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}