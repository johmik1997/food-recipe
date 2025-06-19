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
	"github.com/go-playground/validator/v10"
)
func UpdateRecipe() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var request requests.UpdateRecipeRequest
		log.Println("request", request)
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		validate := validator.New()
		if err := validate.Struct(request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		log.Println("Received update recipe request id: ", request.Input.ID)

		var query struct {
			Recipe []struct {
				ID            graphql.Int    `graphql:"id"`
				Title         graphql.String `graphql:"title"`
				UserId        graphql.Int    `graphql:"user_id"`
				FeaturedImage graphql.String `graphql:"featured_image"`
			} `graphql:"recipes(where: {id: {_eq: $id}})"`
		}

		queryVars := map[string]interface{}{
			"id": graphql.Int(request.Input.ID),
		}

		log.Println("queryVars id", queryVars)

		if err := client.Query(ctx, &query, queryVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query recipe", "details": err.Error()})
			return
		}

		if len(query.Recipe) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
			return
		}

		log.Println("query.Recipe[0].userid", query.Recipe[0].UserId)

		if query.Recipe[0].UserId != graphql.Int(request.Input.UserId) {
			c.JSON(http.StatusForbidden, gin.H{"message": "You are not allowed to update this recipe"})
			return
		}

		imageUrl, exists := c.Get("imageUrl")
		if !exists {
			imageUrl = ""
		}

		featuredImage, ok := imageUrl.(string)
		if !ok || featuredImage == "" {
			featuredImage = string(query.Recipe[0].FeaturedImage)
		}

		var mutation struct {
			UpdateRecipe struct {
				ID              graphql.Int    `graphql:"id"`
				Title           graphql.String `graphql:"title"`
				Description     graphql.String `graphql:"description"`
				PreparationTime graphql.Int    `graphql:"prep_time"`
				CookTime     int                 `graphql:"cook_time"`
				Servings        int                 `graphql:"servings"`
				FeaturedImage   graphql.String `graphql:"featured_image"`
				CategoryId      graphql.Int    `graphql:"category_id"`
				UserId          graphql.Int    `graphql:"user_id"`
				Price           graphql.Int    `graphql:"price"`
			} `graphql:"update_recipes_by_pk(pk_columns: {id: $id}, _set: {title: $title, description: $description, prep_time: $prep_time,cook_time:$cook_time,servings:$servings, featured_image: $featured_image, category_id: $category_id, price: $price, user_id: $user_id})"`
		}

		mutationVars := map[string]interface{}{
			"id":               graphql.Int(request.Input.ID),
			"title":            graphql.String(request.Input.Title),
			"description":      graphql.String(request.Input.Description),
			"prep_time":     graphql.Int(request.Input.PreparationTime),
			"cook_time":      graphql.Int(request.Input.CookTime),
			"servings":       graphql.Int(request.Input.Servings),
			"featured_image":   graphql.String(featuredImage),
			"category_id":      graphql.Int(request.Input.CategoryId),
			"user_id":          graphql.Int(request.Input.UserId),
			"price":            graphql.Int(request.Input.Price),
		}
		log.Println("mutationVarsssss", mutationVars)

		if err := client.Mutate(ctx, &mutation, mutationVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update recipe", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response.UpdateRecipeResponse{Message: "Recipe updated successfully"})
	}
}
