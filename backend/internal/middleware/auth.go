package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"foodrecipe/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	userIDKey contextKey = "user_id"
	tokenKey  contextKey = "token"
)

// (Move this function above AuthMiddleware)

// extractUserIDFromClaims checks multiple possible claim locations for user ID
func extractUserIDFromClaims(claims jwt.MapClaims) string {
	// 1. Check Hasura standard claims first
	if hasuraClaims, ok := claims["https://hasura.io/jwt/claims"].(map[string]interface{}); ok {
		if id, ok := hasuraClaims["x-hasura-user-id"].(string); ok {
			return id
		}
	}

	// 2. Check common JWT claims
	if id, ok := claims["user_id"].(string); ok {
		return id
	}
	
	// 3. Fall back to standard 'sub' claim
	if id, ok := claims["sub"].(string); ok {
		return id
	}

	return ""
}

// AuthMiddleware verifies JWT tokens and injects user ID into context
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// 1. Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "authorization header required")
				return
			}

			// 2. Validate Bearer token format
			if !strings.HasPrefix(authHeader, "Bearer ") {
				utils.WriteJSONError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			// 3. Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// 4. Validate JWT
			claims, err := utils.ValidateJWT(tokenString, []byte(jwtSecret))
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "invalid token: "+err.Error())
				return
			}

			userID := extractUserIDFromClaims(claims)
			if userID == "" {
				utils.WriteJSONError(w, http.StatusUnauthorized,
					"JWT must contain one of: x-hasura-user-id, user_id, or sub claim")
				return
			}

			// ✅ 6. Set both user ID and token in context
			ctx := context.WithValue(r.Context(), tokenKey, tokenString)
			fmt.Printf("Extracted token: %s\n", tokenString)
			ctx = context.WithValue(ctx, userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext safely extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}  

func GetTokenIDFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenKey).(string)
	fmt.Printf("Extracted token from context: %s, ok: %v\n", token, ok)
	return token, ok
}

func AttachTokenToContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}