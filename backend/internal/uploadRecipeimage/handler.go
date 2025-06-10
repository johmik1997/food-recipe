package uploadRecipeimage

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"log"
// 	"net/http"
// )

// type Handler struct {
// 	service *Service
// }

// type HasuraActionRequest struct {
// 	Input struct {
// 		Input UploadRecipeImageInput `json:"input"`
// 	} `json:"input"`
// }



// func NewHandler(service *Service) *Handler {
// 	return &Handler{service: service}
// }

// func (h *Handler) UploadRecipeImageHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
	
// 	var req HasuraActionRequest

// 	// Log raw request body for debugging
// 	bodyBytes, _ := io.ReadAll(r.Body)
// 	log.Printf("Raw request body: %s", string(bodyBytes))
// 	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body for decoding

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		log.Printf("[ImageUpload] Decode error: %v", err)
// 		http.Error(w, `{"message":"failed to parse input"}`, http.StatusBadRequest)
// 		return
// 	}

// 	input := req.Input.Input
// 	log.Printf("[ImageUpload] Received payload: %+v", input)

// 	// Validate required fields
// 	if input.RecipeID == "" {
// 		http.Error(w, `{"message":"recipe_id is required"}`, http.StatusBadRequest)
// 		return
// 	}
// 	if input.Base64Str == "" {
// 		http.Error(w, `{"message":"base64str is required"}`, http.StatusBadRequest)
// 		return
// 	}
// 	if input.Filename == "" {
// 		http.Error(w, `{"message":"filename is required"}`, http.StatusBadRequest)
// 		return
// 	}

// 	imageURL, err := h.service.UploadRecipeImage(r.Context(), input)
// 	if err != nil {
// 		log.Printf("[ImageUpload] Upload error: %v", err)
// 		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"success":   true,
// 		"message":   "Image uploaded and record inserted successfully",
// 		"image_url": imageURL,
// 	}

// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		log.Printf("[ImageUpload] Response encoding error: %v", err)
// 		http.Error(w, `{"message":"failed to encode response"}`, http.StatusInternalServerError)
// 	}
// }