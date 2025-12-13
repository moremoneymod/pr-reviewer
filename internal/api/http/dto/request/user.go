package request

type UserActiveRequest struct {
	UserID   string `json:"user_id" validate:"required,min=1"`
	IsActive bool   `json:"is_active" validate:"required"`
}
