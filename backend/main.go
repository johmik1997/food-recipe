package main

import (
	"foodrecipe/handlers"
	"log"
	"net/http"
)

func main() {
	// Register route handlers
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/upload-recipe-image", handlers.UploadRecipeImageHandler)
	http.HandleFunc("/create-recipe", handlers.CreateRecipeHandler)	
	// Start the server
	log.Println("Server listening on port 3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}