// package imageupload/handler.go
package uploadimage

import (
	"encoding/json"
	"net/http"
	"github.com/cloudinary/cloudinary-go/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type UploadRequest struct {
	Input struct {
		RecipeID   string `json:"recipe_id"`
		Base64Str  string `json:"base64str"`
		Filename   string `json:"filename"`
		IsFeatured bool   `json:"is_featured"`
	} `json:"input"`
}

func (h *Handler) UploadRecipeImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	result, err := h.service.UploadRecipeImage(r.Context(), UploadInput{
		RecipeID:    req.Input.RecipeID,
		Base64Str:   req.Input.Base64Str,
		Filename:    req.Input.Filename,
		IsFeatured:  req.Input.IsFeatured,
	})
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}