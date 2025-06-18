package recipe

import (
	"context"
	"foodrecipe/controllers/libs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)

func GetAllRecipes() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Define the query to fetch all recipes
		var query struct {
			Recipes []struct {
				ID              graphql.Int    `graphql:"id"`
				Title           graphql.String `graphql:"title"`
				Description     graphql.String `graphql:"description"`
				PreparationTime graphql.Int    `graphql:"prep_time"`
				CookingTime     int                 `graphql:"cook_time"`
				Servings        int                 `graphql:"servings"`
				FeaturedImage   graphql.String `graphql:"featured_image"`
				UserId          graphql.Int    `graphql:"user_id"`
				CategoryId      graphql.Int    `graphql:"category_id"`
				Price           graphql.Int    `graphql:"price"`
			} `graphql:"recipes"`
		}

		// Execute the query
		if err := client.Query(ctx, &query, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch recipes", "details": err.Error()})
			return
		}

		// Return the list of recipes
		c.JSON(http.StatusOK, query.Recipes)
	}
}