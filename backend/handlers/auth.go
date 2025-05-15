package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"foodrecipe/models"
	"foodrecipe/utils"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
var jwtSecret = os.Getenv("JWT_SECRET")

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    // Log the entire request body for debugging
    bodyBytes, _ := io.ReadAll(r.Body)
    log.Printf("Raw request body: %s", string(bodyBytes))
    r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
    // Parse input
    var payload struct {
        Input struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        } `json:"input"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, `{"message":"invalid input"}`, http.StatusBadRequest)
        return
    }

    input := payload.Input

    // 1. Find user by email
    query := `
    query ($email: String!) {
      users(where: {email: {_eq: $email}}, limit: 1) {
        id
        email
        password
      }
    }`

    vars := map[string]interface{}{"email": input.Email}
    reqBody, _ := json.Marshal(map[string]interface{}{"query": query, "variables": vars})

    req, _ := http.NewRequest("POST", utils.GraphQLURL, bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-hasura-admin-secret", utils.AdminSecret)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        http.Error(w, `{"message":"database request failed"}`, http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // 2. Handle user lookup
    var result struct {
        Data struct {
            Users []struct {
                ID       string `json:"id"`
                Email    string `json:"email"`
                Password string `json:"password"`
            } `json:"users"`
        } `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.Data.Users) == 0 {
        http.Error(w, `{"message":"invalid credentials"}`, http.StatusUnauthorized)
        return
    }

    user := result.Data.Users[0]

    // 3. Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        http.Error(w, `{"message":"invalid credentials"}`, http.StatusUnauthorized)
        return
    }

    // 4. Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":   user.ID,
        "email": user.Email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
        "https://hasura.io/jwt/claims": map[string]interface{}{
            "x-hasura-allowed-roles": []string{"user"},
            "x-hasura-default-role":  "user",
            "x-hasura-user-id":      user.ID,
        },
    })

    tokenString, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        http.Error(w, `{"message":"token generation failed"}`, http.StatusInternalServerError)
        return
    }

    // 5. Return response matching your GraphQL schema
    response := map[string]interface{}{
        "token": tokenString,
        "user": map[string]interface{}{
            "id":    user.ID,
            "email": user.Email,
        },
    }

    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, `{"message":"response encoding failed"}`, http.StatusInternalServerError)
    }
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload models.ActionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"message":"invalid input"}`, http.StatusBadRequest)
		return
	}

	input := payload.Input.Input

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"message":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// Insert user into Hasura via GraphQL mutation
	mutation := `
mutation ($email: String!, $name: String!, $password: String!) {
  insert_users_one(object: {
    email: $email,
    name: $name,
    password: $password
  }) {
    id
  }
}`

	vars := map[string]interface{}{
		"email":    input.Email,
		"name":     input.Name,
		"password": string(hashedPassword),
	}

	insertReq := models.HasuraInsertRequest{
		Query:     mutation,
		Variables: vars,
	}

	reqBody, _ := json.Marshal(insertReq)

	req, err := http.NewRequest("POST", utils.GraphQLURL, bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, `{"message":"request creation failed"}`, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", utils.AdminSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, `{"message":"failed to execute request"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		http.Error(w, string(body), http.StatusBadRequest)
		return
	}

	var gqlResp models.InsertUserResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		http.Error(w, `{"message":"invalid graphql response"}`, http.StatusInternalServerError)
		return
	}

	userID := gqlResp.Data.InsertUser.ID

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"email": input.Email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-allowed-roles": []string{"user"},
			"x-hasura-default-role":  "user",
			"x-hasura-role":          "user",
			"x-hasura-user-id":       userID,
		},
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, `{"message":"failed to sign token"}`, http.StatusInternalServerError)
		return
	}

	respObj := models.RegisterResponse{
		Token: tokenString,
	}
	respObj.User.ID = userID
	respObj.User.Email = input.Email

	json.NewEncoder(w).Encode(respObj)
}