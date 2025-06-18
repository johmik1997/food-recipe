package auth
import (
	"context"
	"foodrecipe/controllers/helpers"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
func DeleteUserById() gin.HandlerFunc {
	return func(c *gin.Context) {

		client := libs.SetupGraphqlClient()

		var input requests.DeleteUserWithIdInput
		if err := c.ShouldBindBodyWithJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		var query struct {
			User struct {
				ID       graphql.Int    `graphql:"id"`
				Email    graphql.String `graphql:"email"`
				UserName graphql.String `graphql:"username"`
			} `graphql:"users_by_pk(id: $id)"`
		}

		queryVars := map[string]interface{}{
			"id": graphql.Int(input.UserID),
		}

		if err := client.Query(context.Background(), &query, queryVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query user data", "details": err.Error()})
			return
		}

		var mutation struct {
			DeleteUser struct {
				ID graphql.Int `graphql:"id"`
			} `graphql:"delete_users_by_pk(id: $id)"`
		}

		mutationVars := map[string]interface{}{
			"id": graphql.Int(input.UserID),
		}

		if err := client.Mutate(context.Background(), &mutation, mutationVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user", "details": err.Error()})
			return
		}
		emailData := helpers.EmailData{
			Name:    string(query.User.UserName),
			Email:   string(query.User.Email),
			Subject: "Account Deletion Confirmation",
		}

		success, errorString := helpers.SendEmail(
			[]string{emailData.Email},
			"deletedAccount.html", emailData,
		)

		if !success {
			// Undo the deletion if email fails
			var undoMutation struct {
				InsertUser struct {
					ID graphql.Int `graphql:"id"`
				} `graphql:"insert_users_one(object: {id: $id, email: $email, user_name: $userName})"`
			}
			undoVars := map[string]interface{}{
				"id":       graphql.Int(query.User.ID),
				"email":    graphql.String(query.User.Email),
				"userName": graphql.String(query.User.UserName),
			}
			if undoErr := client.Mutate(context.Background(), &undoMutation, undoVars); undoErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to rollback deletion", "details": undoErr.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email, deletion rolled back", "details": errorString})
			return
		}

		c.JSON(http.StatusOK, response.DeleteUserWithEmailResponse{
			Status:  "success",
			Message: "User deleted and email sent successfully",
		})
	}
}