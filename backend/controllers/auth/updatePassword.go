// Reset Password godoc
// @Summary Reset password
// @Description Send a password reset token to the user's email address if they have a verified email.
// @Tags Reset Password
// @Accept json
// @Produce json
// @Param request body requests.UpdatePasswordRequest true "Password Reset Request"
// @Success 200 {object} response.UpdatePasswordResponse "Password reset request success"
// @Failure 400 {object} gin.H "Invalid input"
// @Failure 401 {object} gin.H "Unauthorized - Invalid credentials or unverified email"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/update-password [post]
package auth
import (
	"context"
	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
func UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		var request requests.UpdatePasswordRequest

		if err := c.ShouldBind(&request); err != nil {
			log.Printf("error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input", "details": err.Error()})
			return
		}
		log.Printf("Incoming request: %+v", request)
		if validationError := validate.Struct(request); validationError != nil {
			log.Printf("Validation error: %v", validationError)
			c.JSON(http.StatusBadRequest, gin.H{"message": "validation failed", "details": validationError.Error()})
			return
		}

		var query struct {
			Tokens []struct {
				ID     graphql.Int    `graphql:"id"`
				Token  graphql.String `graphql:"token"`
				UserId graphql.Int    `graphql:"user_id"`
			} `graphql:"email_verification_tokens(where: {token: {_eq: $token}})"`
		}
		queryVars := map[string]interface{}{
			"token": graphql.String(request.Input.Token),
		}

		if err := client.Query(context.Background(), &query, queryVars); err != nil {
			log.Printf("failed to query token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to query token", "details": err.Error()})
			return
		}

		if len(query.Tokens) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid tokens"})
			return
		}

		password := helpers.HashPassword(request.Input.Password)

		var mutation struct {
			UpdateUser struct {
				ID       graphql.Int    `graphql:"id"`
				UserName graphql.String `graphql:"username"`
				Email    graphql.String `graphql:"email"`
				Profile  graphql.String `graphql:"profile"`
				Role     graphql.String `graphql:"role"`
			} `graphql:"update_users_by_pk(pk_columns: {id: $id}, _set: {password: $password})"`
		}

		mutationVars := map[string]interface{}{
			"id":       graphql.Int(request.Input.UserId),
			"password": graphql.String(password),
		}

		err := client.Mutate(context.Background(), &mutation, mutationVars)
		if err != nil {
			log.Printf("failed to update the password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update the password", "details": err.Error()})
			return
		}
		emailForm := helpers.EmailData{
			Name:    string(mutation.UpdateUser.UserName),
			Email:   string(mutation.UpdateUser.Email),
			Subject: "Password reseted successfully!",
		}

		sucess, errorString := helpers.SendEmail(
			[]string{emailForm.Email}, "resetPasswordSuccess.html", emailForm,
		)
		if !sucess {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to send password reset email", "details": errorString})
			return
		}

		response := response.UpdatePasswordResponse{
			Message: "your password has been reseted successfully",
		}
		log.Println(response)

		c.JSON(http.StatusOK, response)

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

	}
}