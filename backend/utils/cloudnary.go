package utils

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	graphql "github.com/hasura/go-graphql-client"
)

var (
	CloudName   = os.Getenv("CLOUDINARY_CLOUD_NAME")
	APIKey      = os.Getenv("CLOUDINARY_API_KEY")
	APISecret   = os.Getenv("CLOUDINARY_API_SECRET")
   // e.g., http://graphql-engine:8080/v1/graphql
)

var GraphQLClient *graphql.Client

func InitGraphQLClient() {
	GraphQLClient = graphql.NewClient(
		"http://graphql-engine:8080/v1/graphql",
		&http.Client{},
	).WithRequestModifier(func(req *http.Request) {
		req.Header.Set("x-hasura-admin-secret", "myadminsecretkey")
	})
}
	
	// headerRoundTripper is a custom RoundTripper to add headers to requests
	type headerRoundTripper struct {
		underlyingTransport http.RoundTripper
		headers             map[string]string
	}
	
	func (h *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
		for key, value := range h.headers {
			req.Header.Set(key, value)
		}
		return h.underlyingTransport.RoundTrip(req)
	}




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