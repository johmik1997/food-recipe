package uploadimage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service struct {
	cld    *cloudinary.Cloudinary
	client GraphQLClient
}

func NewService(cld *cloudinary.Cloudinary, client GraphQLClient) *Service {
	return &Service{
		cld:    cld,
		client: client,
	}
}

func (s *Service) UploadProfileImage(input UploadProfileInput) (*UploadImageOutput, error) {
	// Clean base64 prefix
	cleanBase64 := input.Base64Str
	if strings.Contains(input.Base64Str, "base64,") {
		parts := strings.Split(input.Base64Str, "base64,")
		if len(parts) > 1 {
			cleanBase64 = parts[1]
		}
	}

	// Decode base64 image
	imgData, err := base64.StdEncoding.DecodeString(cleanBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 image data: %w", err)
	}

	// Upload to Cloudinary
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadResult, err := s.cld.Upload.Upload(ctx, bytes.NewReader(imgData), uploader.UploadParams{
		PublicID:     "profile/" + input.Filename,
		Folder:       "profile",
		ResourceType: "image",
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	log.Println("UserID from service:", input.UserID)


	if err := s.client.UpdateUserAvatar(input.UserID, uploadResult.SecureURL); err != nil {

		log.Printf("Error updating Hasura for user %s: %v", input.UserID, err)
		return nil, fmt.Errorf("hasura update failed: %w", err)
	}
	return &UploadImageOutput{
		Success:  true,
		Message:  "Image uploaded successfully",
		ImageURL: uploadResult.SecureURL,
	}, nil
}
