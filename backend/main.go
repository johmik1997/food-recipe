package main

import (
	"log"
	"net/http"

	"foodrecipe/handlers"
)

func main() {
	// Register route handlers
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/upload-recipe-image", handlers.UploadRecipeImageHandler)

	// Start the server
	log.Println("Server listening on port 3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}