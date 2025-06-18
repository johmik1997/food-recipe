package recipe

import (
	"context"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)

func DeleteRecipe() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Define a struct to match the request body
		var request requests.DeleteRecipeRequest

		// Bind JSON request to struct
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON request", "details": err.Error()})
			return
		}

		log.Printf("Request payload: %+v", request)

		// Validate ID
		if request.Input.RecipeId <= 0 {
			log.Printf("Invalid recipe ID: %d", request.Input.RecipeId)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid recipe ID"})
			return
		}

		// Query the database for the recipe
		var query struct {
			Recipes []struct {
				ID     graphql.Int `graphql:"id"`
				UserId graphql.Int `graphql:"user_id"`
			} `graphql:"recipes(where: {id: {_eq: $id}})"`
		}

		queryVars := map[string]interface{}{
			"id": graphql.Int(request.Input.RecipeId),
		}

		if err := client.Query(ctx, &query, queryVars); err != nil {
			log.Printf("Failed to query recipe data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query recipe data", "details": err.Error()})
			return
		}

		if len(query.Recipes) == 0 {
			log.Printf("Recipe not found: %d", request.Input.RecipeId)
			c.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
			return
		}

		userID := request.Input.UserId

		// Authorization check
		if query.Recipes[0].UserId != graphql.Int(userID) {
			log.Println("Unauthorized deletion attempt userid", query.Recipes[0].UserId, userID)
			log.Printf("Unauthorized deletion attempt: userID=%d, recipeUserID=%d", userID, query.Recipes[0].UserId)
			c.JSON(http.StatusForbidden, gin.H{"message": "You are not allowed to delete this recipe"})
			return
		}

		// Perform deletion
		var mutation struct {
			DeleteRecipe struct {
				ID graphql.Int `graphql:"id"`
			} `graphql:"delete_recipes_by_pk(id: $id)"`
		}

		mutationVars := map[string]interface{}{
			"id": graphql.Int(request.Input.RecipeId),
		}

		if err := client.Mutate(ctx, &mutation, mutationVars); err != nil {
			log.Printf("Failed to delete recipe: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete recipe", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response.RemoveRecipeOutput{
			Message: "recipe deleted successfully",
		})
	}
}
