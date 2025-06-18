package recipe

import (
	"bytes"
	"context"

	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)

// create recipe
func AddRecipe() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var request requests.AddRecipeRequest
		body, _ := c.GetRawData()
		log.Println("Raw request body:", string(body))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		log.Println("incoming request body", request)
		// Initialize validator
		validate := validator.New()
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request input ", "details": err.Error()})
			return
		}
		if err := validate.Struct(request); err != nil {
			log.Printf("Validation error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Validation failed",
				"details": err.Error(),
			})
			return
		}

		log.Println("Received recipe request: ", request)

		// Check if recipe already exists (same title)
		var query struct {
			Recipe []struct {
				ID    graphql.Int    `graphql:"id"`
				Title graphql.String `graphql:"title"`
			} `graphql:"recipes(where: {title: {_eq: $title}})"`
		}

		queryVars := map[string]interface{}{
			"title": graphql.String(request.Input.Title),
		}

		if err := client.Query(ctx, &query, queryVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query recipes", "details": err.Error()})
			return
		}

		if len(query.Recipe) != 0 {
			log.Println("Recipe already exists with the same title")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Recipe already exists"})
			return
		}

		imageUrl, exists := c.Get("imageUrls")
		if !exists {
			imageUrl = ""
		}

		featuredImage, ok := imageUrl.(string)
		if !ok {
			featuredImage = ""
		}

		// Create new recipe
		var mutation struct {
			CreateRecipe struct {
    ID              graphql.Int    `graphql:"id"`
    Title           graphql.String `graphql:"title"`
    Description     graphql.String `graphql:"description"`
    PreparationTime graphql.Int    `graphql:"prep_time"`
    CookTime        graphql.Int    `graphql:"cook_time"` // ✅ fix this
    Servings        graphql.Int    `graphql:"servings"`  // ✅ fix this
    FeaturedImage   graphql.String `graphql:"featured_image"`
    UserId          graphql.Int    `graphql:"user_id"`
    CategoryId      graphql.Int    `graphql:"category_id"`
    Price           graphql.Int    `graphql:"price"`     // ✅ fix this
}`graphql:"insert_recipes_one(object: {title: $title, description: $description, prep_time: $prep_time,cook_time:$cook_time,servings:$servings, featured_image: $featured_image, category_id: $category_id, price: $price, user_id: $user_id})"`
		}

		mutationVars := map[string]interface{}{
			"title":            graphql.String(request.Input.Title),
			"description":      graphql.String(request.Input.Description),
			"prep_time": graphql.Int(request.Input.PreparationTime),
			"cook_time": graphql.Int(request.Input.CookTime),
			"servings": graphql.Int(request.Input.Servings),
			"featured_image":   graphql.String(featuredImage),
			"category_id":      graphql.Int(request.Input.CategoryId),
			"price":            graphql.Int(request.Input.Price),
			"user_id":          graphql.Int(request.Input.UserId),
		}

		if err := client.Mutate(ctx, &mutation, mutationVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create recipe", "details": err.Error()})
			return
		}

		// Insert ingredients
		for _, ingredient := range request.Input.Ingredients {
			var ingredientMutation struct {
				CreateIngredient struct {
					ID       graphql.Int    `graphql:"id"`
					Name     graphql.String `graphql:"name"`
					Quantity graphql.String `graphql:"quantity"`
					RecipeId graphql.Int    `graphql:"recipe_id"`
					UserId graphql.Int	`graphql:"user_id"`
				} `graphql:"insert_ingredients_one(object: {name: $name, quantity: $quantity, recipe_id: $recipe_id, user_id: $user_id})"`
			}

			ingredientVars := map[string]interface{}{
				"name":      graphql.String(ingredient.Name),
				"quantity":  graphql.String(ingredient.Quantity),
				"recipe_id": graphql.Int(mutation.CreateRecipe.ID),
				"user_id":	graphql.Int(request.Input.UserId),
			}

	log.Printf("Received quantity: %v, type: %T", ingredient.Quantity, ingredient.Quantity)


			if err := client.Mutate(ctx, &ingredientMutation, ingredientVars); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create ingredient", "details": err.Error()})
				return
			}
		}

		// Insert steps
		for _, step := range request.Input.Steps {
			var stepMutation struct {
				CreateStep struct {
					ID          graphql.Int    `graphql:"id"`
					StepNumber  graphql.Int    `graphql:"step_number"`
					Instruction graphql.String `graphql:"instruction"`
					RecipeId    graphql.Int    `graphql:"recipe_id"`
					UserId 		graphql.Int    `graphql:"user_id"`
				} `graphql:"insert_steps_one(object: {step_number: $step_number, instruction: $instruction, recipe_id: $recipe_id, user_id: $user_id})"`
			}

			stepVars := map[string]interface{}{
				"step_number": graphql.Int(step.StepNumber),
				"instruction": graphql.String(step.Instruction),
				"recipe_id":   graphql.Int(mutation.CreateRecipe.ID),
				"user_id":	 graphql.Int(request.Input.UserId),
			}

			if err := client.Mutate(ctx, &stepMutation, stepVars); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create step", "details": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, response.AddRecipeResponseOutput{
			ID:      int(mutation.CreateRecipe.ID),
			Title:   string(mutation.CreateRecipe.Title),
			Message: "Recipe created successfully",
		})
	}
}