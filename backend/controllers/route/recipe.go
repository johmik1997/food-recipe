package route

import (
	"foodrecipe/controllers/middleware"
	"foodrecipe/controllers/recipe"

	"github.com/gin-gonic/gin"
)

func RecipeRoutes(router *gin.Engine) {
	RecipeRoutes := router.Group("/api/recipes")
	{
		RecipeRoutes.POST("/create", middleware.ImageUpload(), recipe.AddRecipe())
		RecipeRoutes.DELETE("/delete", recipe.DeleteRecipe())
		RecipeRoutes.PUT("/update", middleware.ImageUpload(), recipe.UpdateRecipe())
		RecipeRoutes.GET("/", recipe.GetAllRecipes())
		RecipeRoutes.POST("/uploadImg", middleware.ImageUpload(),recipe.UploadImage())
		RecipeRoutes.POST("/updateImg", middleware.ImageUpload(), recipe.UpdateImage())
		RecipeRoutes.POST("/buy-recipe", recipe.BuyRecipe())
		RecipeRoutes.PUT("/verify-payment", recipe.ProcessPayment())
		RecipeRoutes.GET("/callback", recipe.PaymentCallback())
	}
}