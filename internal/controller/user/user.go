package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	contextModel "userwalletservice/internal/model/context"
	userModel "userwalletservice/internal/model/user"
	userService "userwalletservice/internal/service/user"
)

type Controller struct {
	userService *userService.Service
}

func New(userService *userService.Service) *Controller {
	return &Controller{userService: userService}
}

func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req userModel.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := c.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		// 🔥 Ошибка берется из userModel
		if errors.Is(err, userModel.ErrEmailAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict) // 409
			return
		}

		log.Printf("failed to register user: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError) // 500
		return
	}

	resp := userModel.UserResponse{
		ID:    createdUser.ID,
		Email: createdUser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req userModel.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "login - invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	token, err := c.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		// 🔥 Тоже забираем из userModel, если ты её туда перенес
		if errors.Is(err, userModel.ErrInvalidCredentials) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized) // 401
			return
		}

		log.Printf("failed to login user: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"token":  token,
	})
}

func (c *Controller) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(contextModel.UserIDKey).(int)
	if !ok {
		// 🔥 Логируем для себя критическую ошибку конфигурации, чтобы сразу починить
		log.Printf("[ERROR] AuthMiddleware missing or failed to set UserIDKey in context for GetProfile")

		// Отдаем 500, так как это баг нашей разработки
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user, err := c.userService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, userModel.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound) // 404
			return
		}

		log.Printf("failed to get user profile: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError) // 500
		return
	}

	resp := userModel.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)
}
