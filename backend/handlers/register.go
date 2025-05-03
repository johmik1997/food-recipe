package handlers

import (
	"encoding/json"
	"foodrecipe/hasura"
	"foodrecipe/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-hasura-admin-secret")
}

type RegisterHandler struct {
	HasuraClient *hasura.Client
	JWTSecret    string
}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{
		HasuraClient: hasura.NewClient(),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enableCORS(&w)
	
	if r.Method == "OPTIONS" {
		return
	}

	// Verify content type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, `{"message": "Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
		return
	}

	// Parse input
	var requestBody struct {
		Input struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		} `json:"input"`
	}

	// Read and parse the body just once
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict parsing
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, `{"message": "Invalid JSON input"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestBody.Input.Email == "" || requestBody.Input.Name == "" || requestBody.Input.Password == "" {
		http.Error(w, `{"message": "All fields are required"}`, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(requestBody.Input.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Execute Hasura mutation
	variables := map[string]interface{}{
		"name":     requestBody.Input.Name,
		"email":    requestBody.Input.Email,
		"password": hashedPassword,
	}

	result, err := h.HasuraClient.ExecuteRegister(variables, r.Header)
	if err != nil {
		log.Printf("Error executing Hasura mutation: %v", err)
		http.Error(w, `{"message": "Error communicating with Hasura"}`, http.StatusInternalServerError)
		return
	}

	if len(result.Errors) > 0 {
		errorMessage := result.Errors[0].Message
		http.Error(w, `{"message": "`+errorMessage+`"}`, http.StatusBadRequest)
		return
	}

	// Generate JWT token
	userID := result.Data.InsertUsersOne.ID
	token, err := h.generateToken(userID, requestBody.Input.Name)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		http.Error(w, `{"message": "Error generating token"}`, http.StatusInternalServerError)
		return
	}

	// Return the exact structure Hasura expects
	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":    userID,
			"email": requestBody.Input.Email,
			"name":  requestBody.Input.Name,
		},
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"message": "Error creating response"}`, http.StatusInternalServerError)
	}
}

func (h *RegisterHandler) generateToken(userID int, name string) (string, error) {
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