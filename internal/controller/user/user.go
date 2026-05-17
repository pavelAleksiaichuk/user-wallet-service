package user

import (
	"encoding/json"
	"net/http"
	userServ "userwalletservice/internal/service/user"
	contextModel "userwalletservice/internal/model/context"
)

type UserController struct {
	userService *userServ.UserService
}

func New(userService *userServ.UserService) *UserController {
	return &UserController{userService: userService}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := c.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
}

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "login - invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	token, err := c.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"token":  token,
	})
}

func (c *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(contextModel.UserIDKey).(int)
	if !ok {
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	user, err := c.userService.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	resp := UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)
}
