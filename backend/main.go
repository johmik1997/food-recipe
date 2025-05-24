package main

import (
	"context"
	"foodrecipe/internal/auth"
	"foodrecipe/internal/recipe"
	"foodrecipe/internal/server"
	"foodrecipe/internal/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hasura/go-graphql-client"
)

func main() {
	// Load configuration
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize GraphQL client
	gqlClient := graphql.NewClient(
		cfg.HasuraEndpoint,
		http.DefaultClient,
	).WithRequestModifier(func(r *http.Request) {
		r.Header.Set("x-hasura-admin-secret", cfg.HasuraAdminSecret)
	})

	// Auth setup
	userClient := auth.NewUserClient(gqlClient)
	authService := auth.NewService(userClient, cfg)
	authCfg := &auth.Config{
		JWTSecret: cfg.JWTSecret,
	}
	authHandler := auth.NewHandler(authService, authCfg)

	// Recipe setup
	recipeClient := recipe.NewRecipeClient(cfg.HasuraEndpoint, cfg.JWTSecret)
	recipeService := recipe.NewRecipeService(*recipeClient)
	recipeHandler := recipe.NewRecipeHandler(*recipeService)

	// Create and configure server
	srv := server.NewServer(cfg, authHandler, recipeHandler)
	srv.SetupRoutes()

	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on :%s", cfg.HTTPPort)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-shutdownChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server shut down gracefully.")
}
