package utils

import (
	"time"
	"foodrecipe/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-256-bit-secret") // Should be from environment in production

type Claims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"https://hasura.io/jwt/claims": jwt.MapClaims{
			"x-hasura-allowed-roles": []string{"user"},
			"x-hasura-default-role":  "user",
			"x-hasura-user-id":      user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}