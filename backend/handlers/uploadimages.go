package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"foodrecipe/models"
	"foodrecipe/utils"
	"io"
	"log"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadRecipeImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload models.UploadRecipeImageRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Decode error:", err)
		http.Error(w, `{"message":"failed to parse input"}`, http.StatusBadRequest)
		return
	}

	input := payload.Input

	// Check if necessary fields are provided
	if input.Base64Str == "" || input.Filename == "" || input.RecipeID == "" {
		http.Error(w, `{"message":"missing required fields"}`, http.StatusBadRequest)
		return
	}

	// Decode base64
imageData, err := base64.StdEncoding.DecodeString(input.Base64Str)
if err != nil {
	http.Error(w, `{"message":"invalid base64 image"}`, http.StatusBadRequest)
	return
}

// Upload to Cloudinary
cld, err := cloudinary.NewFromParams(utils.CloudName, utils.APIKey, utils.APISecret)
if err != nil {
	log.Printf("Cloudinary upload error: %+v\n", err)
	http.Error(w, `{"message":"cloudinary config error"}`, http.StatusInternalServerError)
	return
}

uploadResp, err := cld.Upload.Upload(context.Background(), bytes.NewReader(imageData), uploader.UploadParams{
	PublicID:       "recipes/" + input.Filename,
	Folder:         "recipes",
	UniqueFilename: func(b bool) *bool { return &b }(false),
	Overwrite:      func(b bool) *bool { return &b }(false),
})
if err != nil {
	log.Println("Cloudinary upload error:", err)
	http.Error(w, `{"message":"image upload failed"}`, http.StatusInternalServerError)
	return
}

imageURL := uploadResp.SecureURL


	// Create GraphQL mutation to insert image record into Hasura
	mutation := `
	mutation insertRecipeImage($recipe_id: uuid!, $image_url: String!, $is_featured: Boolean!) {
	  insert_recipe_images_one(object: {
		recipe_id: $recipe_id, 
		image_url: $image_url, 
		is_featured: $is_featured
	  }) {
		id
		image_url
	  }
	}`
	fmt.Println("recipe_id:", input.RecipeID)	
	fmt.Println("image_url:", imageURL)
	fmt.Println("is_featured:", input.IsFeatured)
	variables := map[string]interface{}{
		"recipe_id":  input.RecipeID,
		"image_url":  imageURL,
		"is_featured": input.IsFeatured,
	}

	// Prepare the request to send to Hasura
	reqBody, _ := json.Marshal(map[string]interface{}{"query": mutation, "variables": variables})
	req, err := http.NewRequest("POST", utils.GraphQLURL, bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, `{"message":"failed to create request"}`, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", utils.AdminSecret)

	// Send the request to Hasura
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, `{"message":"failed to execute request"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Println("Hasura error response:", string(body))
		http.Error(w, `{"message":"Hasura request failed"}`, http.StatusInternalServerError)
		return
	}
	
	// Return a success response
	response := models.UploadResponse{
		Success:  true,
		Message:  "Image uploaded and record inserted successfully",
		ImageURL: imageURL,
	}

	// Encode and send the response back to the client
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"message":"failed to encode response"}`, http.StatusInternalServerError)
	}
}