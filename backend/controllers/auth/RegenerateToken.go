package auth

import (
	"context"
	"fmt"
	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
func RegenerateEmailVerificationToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var req requests.RegenerateEmailVerificationToken
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("Validation error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		if validationErr := validate.Struct(req); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "details": validationErr.Error()})
			return
		}

		// Fetch the user with the provided email
		var userQuery struct {
			User []struct {
				ID    graphql.Int    `graphql:"id"`
				Name  graphql.String `graphql:"username"`
				Email graphql.String `graphql:"email"`
				IsEmailVerified graphql.Boolean `graphql:"is_email_verified "`
			} `graphql:"users(where: {email: {_eq: $email}})"`
		}
		queryVar := map[string]interface{}{
			"email": graphql.String(req.Input.Email),
		}

		if err := client.Query(ctx, &userQuery, queryVar); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query user data", "details": err.Error()})
			return
		}

		// Check if the user is found
		if len(userQuery.User) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			return
		}

		if (userQuery.User[0].IsEmailVerified ){
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Email already verified"})
			log.Println("your email is already verified!")
			return
		}

		// Fetch the existing email verification token
		var tokenQuery struct {
			Tokens []struct {
				ID        graphql.Int    `graphql:"id"`
				Token     graphql.String `graphql:"token"`
				UserId    graphql.Int    `graphql:"user_id"`
				ExpiresAt graphql.String `graphql:"expires_at"`
			} `graphql:"email_verification_tokens(where: {user_id: {_eq: $userId}})"`
		}

		tokenVar := map[string]interface{}{
			"userId": graphql.Int(userQuery.User[0].ID),
		}

		err := client.Query(ctx, &tokenQuery, tokenVar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query the verification token", "details": err.Error()})
			return
		}
		fmt.Println("the fetched token details", tokenQuery)

		// Token generation
		token, err := helpers.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token", "details": err.Error()})
			return
		}

		// Delete the existing token if there is any
		if len(tokenQuery.Tokens) > 0 {
			var deleteMutation struct {
				DeleteToken struct {
					ID graphql.Int `graphql:"id"`
				} `graphql:"delete_email_verification_tokens_by_pk(id: $id)"`
			}

			deleteVar := map[string]interface{}{
				"id": graphql.Int(tokenQuery.Tokens[0].ID),
			}

			err = client.Mutate(ctx, &deleteMutation, deleteVar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete token", "details": err.Error()})
				return
			}
			log.Printf("user %d successfully deleted old token", userQuery.User[0].ID)
		}

		// Insert a new token
		var insertMutation struct {
			InsertToken struct {
				ID graphql.Int `graphql:"id"`
			} `graphql:"insert_email_verification_tokens_one(object: {token: $token, user_id: $userId})"`
		}

		insertVar := map[string]interface{}{
			"token":  graphql.String(token),
			"userId": graphql.Int(userQuery.User[0].ID),
		}

		log.Printf("Insert Variables: %+v", insertVar)

		err = client.Mutate(ctx, &insertMutation, insertVar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to insert token", "details": err.Error()})
			return
		}
		log.Printf("user %d successfully inserted with token %s", userQuery.User[0].ID, token)

		// Send verification email
		if err := sendVerificationEmail(int(userQuery.User[0].ID), string(userQuery.User[0].Name), string(userQuery.User[0].Email), token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email", "details": err.Error()})
			return
		}

		respose := response.RegenerateEmailVerificationTokenResponse{
			Message: graphql.String("Verification token regenerated and email sent"),
		}

		c.JSON(http.StatusOK,respose)
	}
}
