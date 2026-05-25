package dto

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}
