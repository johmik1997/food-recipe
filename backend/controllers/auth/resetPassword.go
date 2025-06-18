// ResetPassword  request godoc
// @Summary sending Reset password request
// @Description Send a password reset token to the user's email address if they have a verified email.
// @Tags Password Reset Request
// @Accept json
// @Produce json
// @Param request body requests.PasswordResetRequest true "Password Reset Request"
// @Success 200 {object} response.ResetRequestOutput "Password reset request success"
// @Failure 400 {object} gin.H "Invalid input"
// @Failure 401 {object} gin.H "Unauthorized - Invalid credentials or unverified email"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/reset-password [post]
package auth

import (
	"context"
	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		var request requests.PasswordResetRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("invalid input:%v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input", "details": err.Error()})
			return
		}
		// fetch user by email
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
		queryVars := map[string]interface{}{
			"email": graphql.String(request.Input.Email),
		}

		if err := client.Query(context.Background(), &query, queryVars); err != nil {
			log.Printf("failed to query a user data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to query user data", "details": err.Error()})
			return
		}

		if len(query.User) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
			return
		}

		user := query.User[0]
		if !user.IsEmailVerified {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please verify your email first to reset your password"})
			return
		}

		token, err := helpers.GenerateToken()
		if err != nil {
			log.Printf("Failed to generate token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token", "details": err.Error()})
			return
		}
		// Insert reset token into the database
		var mutation struct {
			InsertedToken struct {
				ID    graphql.Int    `graphql:"id"`
				Token graphql.String `graphql:"token"`
			} `graphql:"insert_email_verification_tokens_one(object: {token: $token, user_id: $userId})"`
		}

		mutationVars := map[string]interface{}{
			"token":  graphql.String(token),
			"userId": graphql.Int(user.ID),
		}

		if err := client.Mutate(context.Background(), &mutation, mutationVars); err != nil {
			log.Printf("failed to register reset token")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to register reste token", "details": err.Error()})
			return
		}

		// update the user table to add tokenid
		var UpdateUserTokenMutation struct {
			UpdatedUser struct {
				ID      graphql.Int `graphql:"id"`
				TokenID graphql.Int `graphql:"tokenId"`
			} `graphql:"update_users_by_pk(pk_columns: {id: $userId}, _set: {tokenId: $tokenId})"`
		}

		updatedUserVariables := map[string]interface{}{
			"userId":  graphql.Int(user.ID),
			"tokenId": mutation.InsertedToken.ID,
		}

		err = client.Mutate(context.Background(), &UpdateUserTokenMutation, updatedUserVariables)

		if err != nil {
			log.Printf("error updating user with tokenId: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user tokenid", "details": err.Error()})
			return
		}
		log.Printf("User %d successfully updated with tokenId %d", user.ID, UpdateUserTokenMutation.UpdatedUser.TokenID)

		verificationLink := os.Getenv("RESET_PASS_URL") + "/password-reset?token=" + token + "&id=" + strconv.Itoa(int(user.ID))

		// Send password reset email
		emailData := helpers.EmailData{
			Name:    string(user.Name),
			Email:   string(user.Email),
			Link:    verificationLink,
			Subject: "Reset your password",
		}

		if success, errString := helpers.SendEmail(
			[]string{emailData.Email},
			"passReset.html",
			emailData,
		); !success {
			log.Printf("Failed to send password reset email: %v", errString)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send password reset email", "details": errString})
			return
		}

		response := response.ResetRequestOutput{
			ID:      mutation.InsertedToken.ID,
			Message: "we have sent you an email to reset your password, please go and check your email to reset your password",
		}

		c.JSON(http.StatusOK, response)

	}
}
