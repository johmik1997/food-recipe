// package main

// import (
// 	"context"
// 	"foodrecipe/internal/auth"
// 	"foodrecipe/internal/recipe"
// 	"foodrecipe/internal/server"
// 	"foodrecipe/internal/uploadRecipeimage"
// 	"foodrecipe/internal/uploadimage"

// 	"foodrecipe/internal/utils"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/cloudinary/cloudinary-go/v2"
// 	"github.com/hasura/go-graphql-client"
// )

// func main() {
// cfg, err := utils.LoadConfig()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	// Initialize GraphQL client for Hasura
// 	gqlClient := graphql.NewClient(cfg.HasuraEndpoint, http.DefaultClient).
// 		WithRequestModifier(func(r *http.Request) {
// 			r.Header.Set("x-hasura-admin-secret", cfg.HasuraAdminSecret)
// 			if cfg.JWTSecret != "" {
// 				r.Header.Set("Authorization", "Bearer "+cfg.JWTSecret)
// 			}
// 		})

// // Initialize Cloudinary
// 	cld, err := cloudinary.NewFromParams(
// 		cfg.CloudinaryCloudName,
// 		cfg.CloudinaryAPIKey,
// 		cfg.CloudinaryAPISecret,
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize Cloudinary: %v", err)
// 	}

// 	// Auth setup
// 	userClient := auth.NewUserClient(gqlClient, cfg.HasuraEndpoint, cfg.HasuraAdminSecret)
// 	authService := auth.NewService(userClient, cfg)
// 	authHandler := auth.NewHandler(authService, &auth.Config{
// 		JWTSecret: cfg.JWTSecret,
// 	})

// 	// Recipe setup
// // In your main.go:
// recipeClient := recipe.NewClient(gqlClient)
// recipeService := recipe.NewService(recipeClient)
// recipeHandler := recipe.NewHandler(recipeService)

// //http.Handle("/recipes", middleware.AuthMiddleware(http.HandlerFunc(recipeHandler.CreateRecipe)))
// //cld, _ := cloudinary.NewFromParams(utils.CloudName, utils.APIKey, utils.APISecret)
// recipeimageclient := uploadRecipeimage.NewHasuraClient()
// recipeimageservice := uploadRecipeimage.NewService(cld, recipeimageclient)
// recipeimagehandler := uploadRecipeimage.NewHandler(recipeimageservice)

// 	imageClient := uploadimage.NewClient(uploadimage.ClientConfig{
// 		URL:         cfg.HasuraEndpoint,
// 		AdminSecret: cfg.HasuraAdminSecret,
// 	})
// 	imageService := uploadimage.NewService(cld, imageClient)
// 	imageHandler := uploadimage.NewHandler(imageService)

// // Create server instance
// 	srv := server.NewServer(
// 		cfg,
// 		authHandler,
// 		recipeHandler,
// 		imageHandler,
// 		recipeimagehandler,
// 	)
// 	srv.SetupRoutes()

// 	// Setup graceful shutdown
// 	shutdownChan := make(chan os.Signal, 1)
// 	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	// Start server
// 	go func() {
// 		log.Printf("Server starting on port %s", cfg.HTTPPort)
// 		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Server error: %v", err)
// 		}
// 	}()

// 	// Wait for shutdown signal
// 	<-shutdownChan
// 	log.Println("Shutdown signal received")

// 	// Create shutdown context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Attempt graceful shutdown
// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Printf("Graceful shutdown failed: %v", err)
// 	} else {
// 		log.Println("Server stopped gracefully")
// 	}
// }

package main

import (
	"fmt"
	"foodrecipe/controllers/route"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

var store *sessions.CookieStore

// @title User API for  recipes project
// @version 1.0
// @description This is the API documentation for user api.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email liduhon3@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5000
// @BasePath /api
func main() {
	err := godotenv.Load(`../.env`)
	if err != nil {
		fmt.Println("error loading environment variables", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Load session secret from environment variable
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET is not set in environment variables")
		sessionSecret = "secret"
	}

	// Assign the session store to `gothic` using the session secret
	store = sessions.NewCookieStore([]byte(sessionSecret))
	gothic.Store = store

	router := gin.New()
	router.Use(gin.Logger())
	
	route.RegisterRoutes(router)
	fmt.Printf("Server running on port %s", port)

	router.Run(":" + port)
}