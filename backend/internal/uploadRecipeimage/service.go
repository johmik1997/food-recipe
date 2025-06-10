package uploadRecipeimage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Service handles the business logic for recipe image uploads
type Service struct {
	cld    *cloudinary.Cloudinary
	client GraphQLClient
}

// Handler handles HTTP requests for recipe image uploads
type Handler struct {
	service *Service
}

// HasuraActionRequest represents the incoming request structure from Hasura
type HasuraActionRequest struct {
	Input struct {
		Input UploadRecipeImageInput `json:"input"`
	} `json:"input"`
}

// UploadRecipeImageInput contains the data needed to upload a recipe image
type UploadRecipeImageInput struct {
	Base64Str  string `json:"base64Str"`
	Filename   string `json:"filename"`
	RecipeID   string `json:"recipe_id"`
	IsFeatured bool   `json:"is_featured"`
}

// UploadRecipeImageOutput contains the response data after uploading
type UploadRecipeImageOutput struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"`
}

// NewService creates a new Service instance
func NewService(cld *cloudinary.Cloudinary, client GraphQLClient) *Service {
	return &Service{
		cld:    cld,
		client: client,
	}
}

// NewHandler creates a new Handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// UploadRecipeImageHandler handles HTTP requests for uploading recipe images
func (h *Handler) UploadRecipeImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req HasuraActionRequest

	// Log and read request body
	bodyBytes, _ := io.ReadAll(r.Body)
	log.Printf("[RecipeImageUpload] Raw request body: %s", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[RecipeImageUpload] Decode error: %v", err)
		http.Error(w, `{"message":"Invalid request format"}`, http.StatusBadRequest)
		return
	}

	input := req.Input.Input
	log.Printf("[RecipeImageUpload] Received payload: %+v", input)

	// Validate required fields
	if input.RecipeID == "" {
		http.Error(w, `{"message":"recipeId is required"}`, http.StatusBadRequest)
		return
	}
	if input.Base64Str == "" {
		http.Error(w, `{"message":"base64Str is required"}`, http.StatusBadRequest)
		return
	}
	if input.Filename == "" {
		http.Error(w, `{"message":"filename is required"}`, http.StatusBadRequest)
		return
	}

	// Process the upload
	result, err := h.service.UploadRecipeImage(r.Context(), input)
	if err != nil {
		log.Printf("[RecipeImageUpload] Service error: %v", err)
		http.Error(w, fmt.Sprintf(`{"message":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Return successful response
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("[RecipeImageUpload] Response encoding error: %v", err)
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// UploadRecipeImage handles the core logic for uploading recipe images
func (s *Service) UploadRecipeImage(ctx context.Context, input UploadRecipeImageInput) (*UploadRecipeImageOutput, error) {
	log.Printf("[RecipeImageUpload] Starting image upload for recipe: %s", input.RecipeID)
	log.Printf("[RecipeImageUpload] Filename: %s, Featured: %t", input.Filename, input.IsFeatured)

	// Clean base64 prefix if present
	cleanBase64 := strings.TrimPrefix(input.Base64Str, "data:image/")
	if idx := strings.Index(cleanBase64, ";base64,"); idx != -1 {
		cleanBase64 = cleanBase64[idx+8:]
	}

	// Validate input
	if cleanBase64 == "" || input.Filename == "" || input.RecipeID == "" {
		return nil, fmt.Errorf("missing required fields: base64 image, filename, or recipe ID")
	}

	// Decode base64 image
	imgData, err := base64.StdEncoding.DecodeString(cleanBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 image data: %w", err)
	}

	// Upload to Cloudinary
	uploadCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	uploadResult, err := s.cld.Upload.Upload(uploadCtx, bytes.NewReader(imgData), uploader.UploadParams{
		PublicID:       "recipes/" + input.Filename,
		Folder:         "recipes",
		ResourceType:   "image",
		UniqueFilename: func(b bool) *bool { return &b }(false),
		Overwrite:      func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	// Insert into database
	if err := s.client.InsertRecipeImage(ctx, input.RecipeID, uploadResult.SecureURL, input.IsFeatured); err != nil {
		return nil, fmt.Errorf("database operation failed: %w", err)
	}

	// If the image is featured, update the recipe's featured image
	log.Print("featured image from the service .go ",input.IsFeatured)
	if input.IsFeatured {
		if err := s.client.UpdateRecipeFeaturedImage(ctx, input.RecipeID, uploadResult.SecureURL); err != nil {
			return nil, fmt.Errorf("failed to update featured image: %w", err)
		}
	}

	return &UploadRecipeImageOutput{
		Success:  true,
		Message:  "Recipe image uploaded successfully",
		ImageURL: uploadResult.SecureURL,
	}, nil
}