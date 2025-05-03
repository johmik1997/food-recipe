package models
import (
	"time"
	"github.com/golang-jwt/jwt"
)
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type HasuraResponse struct {
	Data struct {
		InsertUsersOne struct {
			ID int `json:"id"`
		} `json:"insert_users_one"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type JWTClaims struct {
	Sub          string `json:"sub"`
	Name         string `json:"name"`
	Iat          int64  `json:"iat"`
	Iss          string `json:"iss"`
	HasuraClaims struct {
		AllowedRoles []string `json:"x-hasura-allowed-roles"`
		UserID       string   `json:"x-hasura-user-id"`
		DefaultRole  string   `json:"x-hasura-default-role"`
		Role         string   `json:"x-hasura-role"`
	} `json:"https://hasura.io/jwt/claims"`
	Exp int64 `json:"exp"`
}

func (j JWTClaims) Valid() error {
	if time.Now().Unix() > j.Exp {
		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}
	return nil
}
// Add these to your existing models
type LoginInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}