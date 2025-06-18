package recipe

import (
	"context"
	"foodrecipe/controllers/libs"
	"foodrecipe/controllers/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)
func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("🔁 UploadImage endpoint hit")

		client := libs.SetupGraphqlClient()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get image URLs from context
		imageUrls, exists := c.Get("imageUrls")
		if !exists {
			log.Println("❌ imageUrls not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Image URLs not found"})
			return
		}
		imageUrlsArray, ok := imageUrls.([]string)
		if !ok || len(imageUrlsArray) == 0 {
			log.Println("❌ Invalid or empty imageUrls array")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid image URLs"})
			return
		}

		featuredIndex, exists := c.Get("featuredImageIndex")
		if !exists {
			log.Println("❌ featuredImageIndex not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "The featuredImageIndex not found"})
			return
		}
		featuredImageIndex, ok := featuredIndex.(int)
		if !ok {
			log.Println("❌ featuredImageIndex is not an int")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid featured image index type"})
			return
		}
		if featuredImageIndex < 0 || featuredImageIndex >= len(imageUrlsArray) {
			log.Printf("❌ Invalid featured image index: %d", featuredImageIndex)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid featured image index"})
			return
		}

		recipeIDRaw, exists := c.Get("recipeId")
		if !exists {
			log.Println("❌ recipeId not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "The recipeId not found"})
			return
		}
		recipeIDInt, ok := recipeIDRaw.(int)
		if !ok {
			log.Println("❌ recipeId is not an int")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid recipeId type"})
			return
		}

		// Update featured image on recipe
		featuredImageUrl := imageUrlsArray[featuredImageIndex]

		var updateRecipeMutation struct {
			UpdateRecipe struct {
				ID graphql.Int `graphql:"id"`
			} `graphql:"update_recipes_by_pk(pk_columns: {id: $recipeId}, _set: {featured_image: $featuredImageUrl})"`
		}

		updateVars := map[string]interface{}{
			"recipeId":         graphql.Int(recipeIDInt),
			"featuredImageUrl": graphql.String(featuredImageUrl),
		}

		if err := client.Mutate(ctx, &updateRecipeMutation, updateVars); err != nil {
			log.Println("❌ Failed to update recipe:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update recipe", "details": err.Error()})
			return
		}

		// Insert all images into recipe_images
		var uploadResponses []response.ImageUploadResponse
		for i, url := range imageUrlsArray {
			isFeatured := (i == featuredImageIndex)

			var insertImageMutation struct {
				InsertRecipeImage struct {
					ID graphql.Int `graphql:"id"`
				} `graphql:"insert_recipe_images_one(object: {image_url: $imageUrl, recipe_id: $recipeId, is_featured: $isFeatured})"`
			}

			insertVars := map[string]interface{}{
				"imageUrl":   graphql.String(url),
				"recipeId":   graphql.Int(recipeIDInt),
				"isFeatured": graphql.Boolean(isFeatured),
			}

			if err := client.Mutate(ctx, &insertImageMutation, insertVars); err != nil {
				log.Printf("❌ Failed to insert image %d: %s", i, err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload image", "details": err.Error()})
				return
			}

			uploadResponses = append(uploadResponses, response.ImageUploadResponse{
				Url: graphql.String(url),
			})
		}

		log.Println("✅ Images uploaded successfully")
		c.JSON(http.StatusOK, uploadResponses)
	}
}
