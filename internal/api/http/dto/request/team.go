package request

type TeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []TeamMemberRequest `json:"members"`
}

type TeamMemberRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
