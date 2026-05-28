package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"pantrypal/backend/internal/platform/auth"
	"pantrypal/backend/internal/repositories"
)

func AuthRequired(tokens *auth.TokenManager, users *repositories.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader || token == "" {
				writeAPIError(w, http.StatusUnauthorized, "unauthorized", "missing bearer token")
				return
			}

			claims, err := tokens.ParseToken(token)
			if err != nil {
				writeAPIError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
				return
			}

			exists, err := users.ExistsByID(r.Context(), claims.UserID)
			if err != nil || !exists {
				writeAPIError(w, http.StatusUnauthorized, "unauthorized", "invalid token user")
				return
			}

			ctx := WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
