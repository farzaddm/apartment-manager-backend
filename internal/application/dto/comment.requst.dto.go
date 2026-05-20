package dto

type CreateCommentRequest struct {
	Body string `json:"body" binding:"required"`
}

type UpdateCommentRequest struct {
	Body string `json:"body"`
}
