package request

import _ "github.com/go-playground/validator/v10"

type TeamRequest struct {
	TeamName string              `json:"team_name" validate:"required"`
	Members  []TeamMemberRequest `json:"members" validate:"required,min=1,dive"`
}

type TeamMemberRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}
