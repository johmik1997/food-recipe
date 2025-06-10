package uploadimage

// import (

// 	"github.com/google/uuid"
// )

type UploadProfileInput struct {
	UserID    string `json:"userId"`
	Base64Str string `json:"base64str"`
	Filename  string `json:"filename"`
}

type UploadImageOutput struct {
	Success        bool   `json:"success"`
	Message       string `json:"message"`
	ImageURL string `json:"avatar_image_url"`
}





