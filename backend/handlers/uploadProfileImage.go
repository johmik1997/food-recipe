package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"foodrecipe/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadProfileInput struct {
	UserID    string `json:"userId"`
	Base64Str string `json:"base64str"`
	Filename  string `json:"filename"`
}

type UploadImageOutput struct {
	Success        bool   `json:"success"`
	Message       string `json:"message"`
	AvatarImageURL string `json:"avatar_image_url"`
}

func UploadProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Read and log the raw request for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Printf("Raw request body: %s", string(bodyBytes))

	// 2. Parse the input in Hasura Action format
	var payload struct {
		Input struct {
			Input UploadProfileInput `json:"input"`
		} `json:"input"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Printf("Failed to parse request: %v", err)
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	input := payload.Input.Input

	// 2. Validate input
	if input.UserID == "" {
		sendErrorResponse(w, "userId is required", http.StatusBadRequest)
		return
	}
	if input.Base64Str == "" {
		sendErrorResponse(w, "base64str is required", http.StatusBadRequest)
		return
	}
	if input.Filename == "" {
		sendErrorResponse(w, "filename is required", http.StatusBadRequest)
		return
	}

	// 3. Process base64 string
	cleanBase64 := input.Base64Str
	if strings.Contains(input.Base64Str, "base64,") {
		parts := strings.Split(input.Base64Str, "base64,")
		if len(parts) > 1 {
			cleanBase64 = parts[1]
		}
	}

	// 4. Upload to Cloudinary
	imgData, err := base64.StdEncoding.DecodeString(cleanBase64)
	if err != nil {
		sendErrorResponse(w, "Invalid base64 image data", http.StatusBadRequest)
		return
	}

	cld, err := cloudinary.NewFromParams(utils.CloudName, utils.APIKey, utils.APISecret)
	if err != nil {
		log.Printf("Cloudinary config error: %v", err)
		sendErrorResponse(w, "Cloudinary configuration failed", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, bytes.NewReader(imgData), uploader.UploadParams{
		PublicID:     "profile/" + input.Filename,
		Folder:       "profile",
		ResourceType: "image",
	})
	if err != nil {
		log.Printf("Cloudinary upload error: %v", err)
		sendErrorResponse(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	// 5. Update Hasura
	if err := updateUserAvatar(input.UserID, uploadResult.SecureURL); err != nil {
		log.Printf("Hasura update error: %v", err)
		sendErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// 6. Return success response
	json.NewEncoder(w).Encode(UploadImageOutput{
		Success:        true,
		Message:        "Profile image updated successfully",
		AvatarImageURL: uploadResult.SecureURL,
	})
}

func sendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"success": false,
	})
}

func updateUserAvatar(userID, avatarURL string) error {
	mutation := `mutation UpdateUserAvatar($user_id: uuid!, $avatar_url: String!) {
		update_users_by_pk(pk_columns: {id: $user_id}, _set: {avatar_image_url: $avatar_url}) {
			id
		}
	}`

	reqBody, _ := json.Marshal(map[string]interface{}{
		"query": mutation,
		"variables": map[string]interface{}{
			"user_id":    userID,
			"avatar_url": avatarURL,
		},
	})

	req, err := http.NewRequest("POST", utils.GraphQLURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", utils.AdminSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hasura returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}