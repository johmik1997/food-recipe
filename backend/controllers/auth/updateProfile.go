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
 
func UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		var req requests.UpdateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("error from jsonbind %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input", "details": err.Error()})
			return
		}

		if validationErr := validate.Struct(req); validationErr != nil {
			fmt.Println("Validation Error:", validationErr.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"updateProfile": validationErr.Error(),
			},
			)
			return
		}

		imageUrl, exists := c.Get("imageUrl")
		if !exists {
			imageUrl = ""
		}

		proPicUrl, ok := imageUrl.(string)
		if !ok {
			proPicUrl = ""
		}

		roleToUse := req.Input.Role
		if roleToUse == "" {
			var query struct {
				Users []struct {
					Role *string `graphql:"role"`
				} `graphql:"users(where: {id: {_eq: $userId}})"`
			}

			queryVars := map[string]interface{}{
				"userId": graphql.Int(req.Input.UserId),
			}

			err := client.Query(context.Background(), &query, queryVars)
			if err != nil {
				log.Println("Error fetching existing role:", err)
			}

			log.Println("Query response to check role:", query.Users)

			if len(query.Users) > 0 && query.Users[0].Role != nil {
				roleToUse = *query.Users[0].Role
				log.Println("Fetched Role from DB:", roleToUse)
			} else {
				roleToUse = "user"
			}
		}

		log.Println("Final Role to Use:", roleToUse)

		if proPicUrl == "" {
			var mutation struct {
				UpdateProfile struct {
					ID graphql.Int `graphql:"id"`
				} `graphql:"update_users_by_pk(pk_columns: {id: $userId}, _set: {username: $userName, phone: $Phone, role: $Role})"`
			}

			mutationVars := map[string]interface{}{
				"userId":   graphql.Int(req.Input.UserId),
				"userName": graphql.String(req.Input.UserName),
				"Phone":    graphql.String(req.Input.Phone),
				"Role":     graphql.String(roleToUse),
			}

			err := client.Mutate(context.Background(), &mutation, mutationVars)
			if err != nil {
				log.Println("Failed to update user profile:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user profile", "details": err.Error()})
				return
			}

			res := response.UpdateResponce{
				Message: "Profile updated successfully",
			}
			c.JSON(http.StatusOK, res)
		} else {
			var mutation struct {
				UpdateProfile struct {
					ID graphql.Int `graphql:"id"`
				} `graphql:"update_users_by_pk(pk_columns: {id: userId}, _set: {username: $userName, phone: $Phone, profile: $Profile, role: $Role})"`
			}

			mutationVars := map[string]interface{}{
				"userId":   graphql.Int(req.Input.UserId),
				"userName": graphql.String(req.Input.UserName),
				"Phone":    graphql.String(req.Input.Phone),
				"Profile":  graphql.String(proPicUrl),
				"Role":     graphql.String(req.Input.Role),
			}

			err := client.Mutate(context.Background(), &mutation, mutationVars)
			if err != nil {
				log.Println("Failed to update user profile:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user profile"})
				return
			}
			res := response.UpdateResponce{
				Message: "Profile updated successfully",
			}

			c.JSON(http.StatusOK, res)
		}

	}

}
