package router

import (
	"context"
	"net/http"
	"strings"
	"userwalletservice/internal/infrastructure"
	"userwalletservice/internal/model"
)

type AuthMiddleware struct {
	jwtManager *infrastructure.JWTManager
}

func NewAuthMiddleware(jwtManager *infrastructure.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

// Handler — это сам мидлвар
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Вытаскиваем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// 2. Проверяем формат "Bearer <токен>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 3. Валидируем токен
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Записываем UserID в контекст запроса, чтобы контроллеры кошелька знали, кто делает запрос
		ctx := context.WithValue(r.Context(), model.UserIDKey, claims.UserID)
		
		// 5. Передаем запрос дальше по цепочке
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
