package models

type UploadInput struct {
	RecipeID   string `json:"recipe_id"`
	IsFeatured bool   `json:"is_featured"`
	Base64Str  string `json:"base64str"`
	Filename   string `json:"filename"`
}

type UploadRecipeImageRequest struct {
	Input UploadInput `json:"input"`
}

type UploadProfileImageInput struct {
	UserID    string `json:"userId"`
	Base64Str string `json:"base64str"`
	Filename  string `json:"filename"`
}

type UploadProfileImageResponse struct {
	Success        bool   `json:"success"`
	Message       string `json:"message"`
	AvatarImageURL string `json:"avatar_image_url"`
}


type UploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url,omitempty"`
}
type UploadprofileResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	AvaterImageURL string `json:"avatar_image_url,omitempty"`
}
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Input struct {
		Input LoginInput `json:"input"`
	} `json:"input"`
}

type GetUserResponse struct {
	Data struct {
		User []struct {
			ID       string `json:"id"`
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	} `json:"data"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}
type RegisterInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ActionPayload struct {
	Input struct {
		Input RegisterInput `json:"input"`
	} `json:"input"`
}

type InsertUserResponse struct {
	Data struct {
		InsertUser struct {
			ID string `json:"id"`
		} `json:"insert_users_one"`
	} `json:"data"`
}

type HasuraInsertRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type RegisterResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}