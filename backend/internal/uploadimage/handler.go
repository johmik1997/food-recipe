package uploadimage

import (
	"bytes"
	"encoding/json"
	"foodrecipe/internal/utils"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}
type HasuraActionRequest struct {
	Input struct {
		Input UploadProfileImageInput `json:"input"`
	} `json:"input"`
}



func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}
type UploadProfileImageRequest struct {
	Input UploadProfileImageInput `json:"input"`
}

type UploadProfileImageInput struct {
	UserID    string `json:"userId"`
	Base64Str string `json:"base64str"`
	Filename  string `json:"filename"`
}
func (h *Handler) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	var req HasuraActionRequest

	bodyBytes, _ := io.ReadAll(r.Body)
	log.Printf("Raw request body: %s", string(bodyBytes))
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	input := req.Input.Input
	log.Printf("Received UploadProfileImage input: %+v", input)

	if input.UserID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "userId is required")
		return
	}
	if input.Base64Str == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "base64str is required")
		return
	}
	if input.Filename == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "filename is required")
		return
	}

	uploadInput := UploadProfileInput{
		UserID:    input.UserID,
		Base64Str: input.Base64Str,
		Filename:  input.Filename,
	}

	result, err := h.service.UploadProfileImage(uploadInput)
	if err != nil {
		log.Printf("Upload profile image error: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to upload profile image")
		return
	}

response := map[string]interface{}{
	"success":          result.Success,
	"message":          result.Message,
	"avatar_image_url": result.ImageURL,
}


	utils.WriteJSON(w, http.StatusOK, response)
}


