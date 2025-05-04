package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"time"
// 	"net/http"
// 	"os"
// 	"io"
	
// )

// type UploadRecipeImageInput struct {
// 	RecipeID   string `json:"recipe_id"`
// 	ImageURL   string `json:"image_url"`    // unused here, kept for compatibility
// 	IsFeatured bool   `json:"is_featured"`
// 	Base64Str  string `json:"base64str"`    // add to support raw image upload
// 	Filename   string `json:"filename"`
// }

// type UploadRecipeImagePayload struct {
// 	Input struct {
// 		Input UploadRecipeImageInput `json:"input"`
// 	} `json:"input"`
// }

// type UploadRecipeImageOutput struct {
// 	Success   bool   `json:"success"`
// 	Message   string `json:"message,omitempty"`
// 	ImageURL  string `json:"image_url,omitempty"`
// }
// func uploadRecipeImageHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	// Parse form with file max size
// 	err := r.ParseMultipartForm(10 << 20) // 10 MB
// 	if err != nil {
// 		http.Error(w, `{"message":"failed to parse form"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Get the file from form
// 	file, handler, err := r.FormFile("image")
// 	if err != nil {
// 		http.Error(w, `{"message":"failed to read image file"}`, http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Save file to local disk (e.g., ./uploads/)
// 	os.MkdirAll("./uploads", os.ModePerm)
// 	filename := time.Now().Format("20060102150405") + "_" + handler.Filename
// 	filePath := "./uploads/" + filename
// 	out, err := os.Create(filePath)
// 	if err != nil {
// 		http.Error(w, `{"message":"failed to save image"}`, http.StatusInternalServerError)
// 		return
// 	}
// 	defer out.Close()
// 	io.Copy(out, file)

// 	// Optional: get recipe_id from form data if you associate image with a recipe
// 	recipeID := r.FormValue("recipe_id") // if needed

// 	// Insert metadata into Hasura
// 	mutation := `
// 	mutation ($url: String!, $recipe_id: uuid) {
// 		insert_recipe_images_one(object: {
// 			image_url: $url,
// 			recipe_id: $recipe_id
// 		}) {
// 			id
// 			image_url
// 		}
// 	}`

// 	imageURL := "/uploads/" + filename // accessible via frontend if served statically

// 	vars := map[string]interface{}{
// 		"url":       imageURL,
// 		"recipe_id": recipeID, // if blank, it can be omitted or set nullable in Hasura
// 	}

// 	reqBody, _ := json.Marshal(map[string]interface{}{
// 		"query":     mutation,
// 		"variables": vars,
// 	})

// 	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		http.Error(w, `{"message":"request creation failed"}`, http.StatusInternalServerError)
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("x-hasura-admin-secret", adminSecret)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		http.Error(w, `{"message":"failed to execute mutation"}`, http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	if resp.StatusCode != 200 {
// 		http.Error(w, string(body), http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write(body)
// }
