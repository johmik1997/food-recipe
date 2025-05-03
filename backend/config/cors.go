package config

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsMiddleware := cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "x-hasura-admin-secret"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		})
		corsMiddleware(c)
		
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}