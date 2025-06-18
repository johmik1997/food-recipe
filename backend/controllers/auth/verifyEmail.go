// VerifyEmail godoc
// @Summary Verify user email
// @Description Verifies a user's email by checking the verification token and updating the user's status.
// @Tags User Verification
// @Accept json
// @Produce json
// @Param request body requests.EmailVerifyRequest true "User email verification request"
// @Success 200 {object} response.VerifyEmailResponse "Email verified successfully"
// @Failure 400 {object} gin.H "Invalid input data"
// @Failure 401 {object} gin.H "Invalid or expired verification token"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/verify-email [post]
package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"

	"github.com/shurcooL/graphql"
)

func VerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		var request requests.EmailVerifyRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input", "details": err.Error()})
			return
		}
		// Validate the request body
		validationError := validate.Struct(request)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": validationError.Error()})
			return

		}

		// Fetch the token from the database based on the provided
		var query struct {
			Tokens []struct {
				ID        graphql.Int    `graphql:"id"`
				Token     graphql.String `graphql:"token"`
				UserId    graphql.Int    `graphql:"user_id"`
				ExpiresAt graphql.String `graphql:"expires_at"`
			} `graphql:"email_verification_tokens(where: {token: {_eq: $token}})"`
		}
		variables := map[string]interface{}{
			"token": graphql.String(request.Input.VerificationToken),
		}
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to query the verification token", "details": err.Error()})
			return
		}
		// Check if token is found
		if len(query.Tokens) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid verification token"})
			return
		}

		//  validate if the token is expired
		expirationTime, err := time.Parse(time.RFC3339, string(query.Tokens[0].ExpiresAt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to parse token expiration date", "details": err.Error()})
			return
		}

		if time.Now().After(expirationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "expired verification token"})
			return
		}

		var mutation struct {
			UpdateUser struct {
				ID         graphql.ID      `graphql:"id"`
				UserName   graphql.String  `graphql:"username"`
				Email      graphql.String  `graphql:"email"`
				Profile    graphql.String  `graphql:"profile"`
				Role       graphql.String  `graphql:"role"`
				IsVerified graphql.Boolean `graphql:"is_email_verified"`
			} `graphql:"update_users_by_pk(pk_columns: {id: $id}, _set: {is_email_verified: $status})"`
		}

		mutationVariables := map[string]interface{}{
			"id":     graphql.Int(request.Input.UserId),
			"status": graphql.Boolean(true),
		}

		err = client.Mutate(context.Background(), &mutation, mutationVariables)
		if err != nil {
			log.Printf("Error updating user email verification status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user email verification status", "details": err.Error()})
			return
		}

		// delete the token after it has been used
		var deleteMutation struct {
			DeleteToken struct {
				ID graphql.Int `graphql:"id"`
			} `graphql:"delete_email_verification_tokens_by_pk(id: $id)"`
		}
		deleteVariables := map[string]interface{}{
			"id": query.Tokens[0].ID,
		}
		err = client.Mutate(context.Background(), &deleteMutation, deleteVariables)
		if err != nil {
			log.Printf("Error deleting email verification token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete email verification token", "details": err.Error()})
			return
		}

		// return success message

		res := response.VerifyEmailResponse{
			Message: string("Email verified successfully"),
		}

		c.JSON(http.StatusOK, res)
	}
}
