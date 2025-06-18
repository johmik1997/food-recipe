// Login handles user authentication
// @Summary User login
// @Description Authenticates a user with email and password, returning access and refresh tokens
// @Tags User Login
// @Accept json
// @Produce json
// @Param request body requests.LoginRequest true "Login request body"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} gin.H "Invalid input or bad request"
// @Failure 401 {object} gin.H "Invalid credentials or unverified email"
// @Failure 500 {object} gin.H "Failed to query user or generate tokens"
// @Router /users/login [post]
package auth
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		var request requests.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		var query struct {
			User []struct {
				ID              graphql.Int     `graphql:"id"`
				Name            graphql.String  `graphql:"username"`
				Email           graphql.String  `graphql:"email"`
				Password        graphql.String  `graphql:"password"`
				Role            graphql.String  `graphql:"role"`
				TokenId         graphql.Int     `graphql:"tokenId"`
				IsEmailVerified graphql.Boolean `graphql:"is_email_verified"`
			} `graphql:"users(where: {email: {_eq: $email}})"`
		}

		variables := map[string]interface{}{
			"email": graphql.String(request.Input.Email),
		}

		if err := client.Query(context.Background(), &query, variables); err != nil {
			log.Printf("failed to query the user with email %s: %v", request.Input.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query user"})
			return
		}

		if len(query.User) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			return

		}
		user := query.User[0]

		// if !user.IsEmailVerified {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "You need to verify your account first to login"})
		// 	return
		// }

		if valid, msg := helpers.VerifyPassword(request.Input.Password, string(user.Password)); !valid {
			c.JSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}

		token, refreshToken, err := helpers.GenerateAllTokens(string(user.Email), string(user.Name), string(user.Role), fmt.Sprintf("%d", user.ID), int(user.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate tokens",
				"details": err.Error(),
			})
			return
		}

		response := response.LoginResponse{
			User: &response.UserResponse{
				ID:           user.ID,
				Token:        graphql.String(token),
				RefreshToken: graphql.String(refreshToken),
				Email:        graphql.String(user.Email),
				Name:         graphql.String(user.Name),
				Role:         graphql.String(user.Role),
			},
		}

		c.JSON(http.StatusOK, response)

	}
}
