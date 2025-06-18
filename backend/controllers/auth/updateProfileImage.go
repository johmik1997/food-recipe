package auth

import (
	"context"
	"fmt"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
func UpdateProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		var req requests.UpdateProfileImage

		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("Error from JSON bind: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		imageUrl, exists := c.Get("imageUrl")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Profile picture is required"})
			return
		}

		proPicUrl, ok := imageUrl.(string)
		if !ok || proPicUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid profile picture URL"})
			return
		}

		// Define the GraphQL mutation
		var mutation struct {
			UpdateProfilePicture struct {
				ID graphql.Int `graphql:"id"`
			} `graphql:"update_users_by_pk(pk_columns: {id: $userId}, _set: {profile: $Profile})"`
		}

		// Mutation variables
		mutationVars := map[string]interface{}{
			"userId":  graphql.Int(req.Input.UserId),
			"Profile": graphql.String(proPicUrl),
		}

		err := client.Mutate(context.Background(), &mutation, mutationVars)
		if err != nil {
			log.Println("Failed to update profile picture:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update profile picture", "details": err.Error()})
			return
		}
		res := response.UpdateProfileResponce{
			Message: "Profile picture updated successfully",
		}
		c.JSON(http.StatusOK, res)
	}
}
