package recipe

import (
	"context"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/requests"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shurcooL/graphql"
)
func UpdateImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := libs.SetupGraphqlClient()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var request requests.ImageUpLoadRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": err.Error()})
			return
		}

		validate := validator.New()
		validationError := validate.Struct(request)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "details": validationError.Error()})
			return
		}

		imageUrls, exists := c.Get("imageUrls")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Image URLs not found"})
			return
		}

		imageUrlsArray, ok := imageUrls.([]string)
		if !ok || len(imageUrlsArray) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid image URLs"})
			return
		}

		if len(imageUrlsArray) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No image URLs provided"})
			return
		}

		var mutation struct {
			DeleteRecipeImage struct {
				AffectedRows graphql.Int `graphql:"affected_rows"`
			} `graohql:"delete_recipe_images(where: {id: {_eq: $recipeImageId}})"`
		}
		MutationVars := map[string]interface{}{
			"recipeImageId": graphql.Int(request.Input.RecipeId),
		}

		err := client.Mutate(ctx, &mutation, MutationVars)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete image", "details": err.Error()})
			return
		}

		var uploadResponse []response.ImageUploadResponse

		for _, imageUrl := range imageUrlsArray {

			var mutation struct {
				InsertRecipeImage struct {
					ID graphql.Int `graphql:"id"`
				} `graphql:"insert_recipe_images_one(object: {image_url: $imageUrl, recipe_id: $recipeId})"`
			}
			mutationVars := map[string]interface{}{
				"imageUrl": graphql.String(imageUrl),
				"recipeId": graphql.Int(request.Input.RecipeId),
			}
			if err := client.Mutate(ctx, &mutation, mutationVars); err != nil {
				log.Println("Error inserting an image:", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload image", "details": err.Error()})
				return
			}

			uploadResponse = append(uploadResponse, response.ImageUploadResponse{
				Url: graphql.String(strconv.Itoa(int(mutation.InsertRecipeImage.ID))),
			})
		}
		c.JSON(http.StatusOK, uploadResponse)

	}
}