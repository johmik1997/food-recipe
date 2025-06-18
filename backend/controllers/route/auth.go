package route

import (
	"foodrecipe/controllers/auth"
	"foodrecipe/controllers/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	userRoutes := incomingRoutes.Group("/api/users")
	{
		userRoutes.POST("/register", middleware.ImageUpload(), auth.RegisterUser())
		userRoutes.POST("/verify-email", auth.VerifyEmail())
		userRoutes.POST("/login", auth.Login())
		userRoutes.POST("/reset-password", auth.ResetPassword())
		userRoutes.POST("/update-password", auth.UpdatePassword())
		userRoutes.POST("/delete", auth.DeleteUser())
		userRoutes.PUT("/update-profile", middleware.ImageUpload(), auth.UpdateProfile())
		userRoutes.DELETE("/delete", auth.DeleteUserById())
		userRoutes.POST("/user", auth.GetUserById())
		userRoutes.POST("/get-token", auth.RegenerateEmailVerificationToken())

	}
}