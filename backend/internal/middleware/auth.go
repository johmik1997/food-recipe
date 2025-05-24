package middleware

import (
	"context"
	"net/http"
	"foodrecipe/internal/utils"
	"strings"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateJWT(tokenString, []byte(jwtSecret))
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), "user_claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}