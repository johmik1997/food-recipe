package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	"foodrecipe/hasura"
	"foodrecipe/models"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	HasuraClient *hasura.Client
	JWTSecret    string
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		HasuraClient: hasura.NewClient(),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse input
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"message": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := h.getUserByEmail(input.Email)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		http.Error(w, `{"message": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, `{"message": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Name)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		http.Error(w, `{"message": "Error generating token"}`, http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := models.LoginResponse{
		Token: token,
		User: models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"login": response,
	})
}

func (h *LoginHandler) getUserByEmail(email string) (*struct {
	ID       int
	Name     string
	Email    string
	Password string
}, error) {
	query := `query GetUserByEmail($email: String!) {
		users(where: {email: {_eq: $email}}, limit: 1) {
			id
			name
			email
			password
		}
	}`

	variables := map[string]interface{}{
		"email": email,
	}

	result, err := h.HasuraClient.ExecuteQuery(query, variables)
	if err != nil {
		return nil, err
	}

	// Parse result and return user
	// Implementation depends on your Hasura response structure
	// This is a placeholder - adjust according to your actual response
	if len(result.Data.Users) == 0 {
		return nil, errors.New("user not found")
	}

	return &result.Data.Users[0], nil
}

func (h *LoginHandler) generateToken(userID int, name string) (string, error) {
	claims := models.JWTClaims{
		Sub:  string(rune(userID)),
		Name: name,
		Iat:  time.Now().Unix(),
		Iss:  "https://myapp.com/",
		HasuraClaims: struct {
			AllowedRoles []string `json:"x-hasura-allowed-roles"`
			UserID       string   `json:"x-hasura-user-id"`
			DefaultRole  string   `json:"x-hasura-default-role"`
			Role         string   `json:"x-hasura-role"`
		}{
			AllowedRoles: []string{"user"},
			UserID:       string(rune(userID)),
			DefaultRole:  "user",
			Role:         "user",
		},
		Exp: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.JWTSecret))
}