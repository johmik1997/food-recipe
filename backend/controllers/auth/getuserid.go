package auth

import (
	"context"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
	
func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		var req requests.GetUserByIdInput
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		var query struct {
			User struct {
				ID       graphql.Int    `graphql:"id"`
				UserName graphql.String `graphql:"username"`
				Email    graphql.String `graphql:"email"`
			} `graphql:"users_by_pk(id: $id)"`
		}

		queryVars := map[string]interface{}{
			"id": graphql.Int(req.Input.UserID),
		}

		if err := client.Query(context.Background(), &query, queryVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query user data", "details": err.Error()})
			return
		}

		response := response.SingleUserResponse{
			ID:       int(query.User.ID),
			UserName: string(query.User.UserName),
			Email:    string(query.User.Email),
		}
		c.JSON(http.StatusOK, response)
	}
}
