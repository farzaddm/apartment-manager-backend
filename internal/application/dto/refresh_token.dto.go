package dto

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshOutput struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
	Message      string      `json:"message"`
}
