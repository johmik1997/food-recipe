package utils

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var (
	CloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
	APIKey    = os.Getenv("CLOUDINARY_API_KEY")
	APISecret = os.Getenv("CLOUDINARY_API_SECRET")
)

func UploadImage(ctx context.Context, imageData []byte, filename string) (string, error) {
	cld, err := cloudinary.NewFromParams(CloudName, APIKey, APISecret)
	if err != nil {
		log.Printf("Cloudinary initialization error: %v", err)
		return "", err
	}

	uploadResp, err := cld.Upload.Upload(ctx, bytes.NewReader(imageData), uploader.UploadParams{
		PublicID:       "recipes/" + filename,
		Folder:         "recipes",
		UniqueFilename: func(b bool) *bool { return &b }(false),
		Overwrite:      func(b bool) *bool { return &b }(false),
	})
	if err != nil {
		log.Printf("Cloudinary upload error: %v", err)
		return "", err
	}

	return uploadResp.SecureURL, nil
}