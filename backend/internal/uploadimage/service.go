// package imageupload/service.go
package uploadimage

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service struct {
	cloudinary *cloudinary.Cloudinary
	gqlClient  *graphqlclient.Client
}

func NewService(cld *cloudinary.Cloudinary, gqlClient *graphqlclient.Client) *Service {
	return &Service{
		cloudinary: cld,
		gqlClient:  gqlClient,
	}
}

type UploadInput struct {
	RecipeID    string
	Base64Str   string
	Filename    string
	IsFeatured  bool
}

type UploadResult struct {
	Success  bool
	Message  string
	ImageURL string
}

func (s *Service) UploadRecipeImage(ctx context.Context, input UploadInput) (*UploadResult, error) {
	// Validate input
	if input.RecipeID == "" || input.Base64Str == "" || input.Filename == "" {
		return nil, errors.New("missing required fields")
	}

	// Decode base64 image
	imageData, err := base64.StdEncoding.DecodeString(
		strings.TrimPrefix(input.Base64Str, "data:image/jpeg;base64,"),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 image: %w", err)
	}

	// Upload to Cloudinary
	uploadResp, err := s.cloudinary.Upload.Upload(ctx, bytes.NewReader(imageData), uploader.UploadParams{
		PublicID:       "recipes/" + input.Filename,
		Folder:         "recipes",
		UniqueFilename: boolPtr(false),
		Overwrite:      boolPtr(false),
	})
	if err != nil {
		return nil, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	
	// Insert into Hasura
	var mutation struct {
		InsertRecipeImage struct {
			ID       string `graphql:"id"`
			ImageURL string `graphql:"image_url"`
		} `graphql:"insert_recipe_images_one(object: {recipe_id: $recipe_id, image_url: $image_url, is_featured: $is_featured})"`
	}

	err = s.gqlClient.Mutate(ctx, &mutation, map[string]interface{}{
		"recipe_id":   input.RecipeID,
		"image_url":   uploadResp.SecureURL,
		"is_featured": input.IsFeatured,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert image record: %w", err)
	}

	return &UploadResult{
		Success:  true,
		Message:  "Image uploaded and record inserted successfully",
		ImageURL: uploadResp.SecureURL,
	}, nil
}

func boolPtr(b bool) *bool { return &b }