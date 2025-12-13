package request

type PRCreateRequest struct {
	PullRequestID   string `json:"pull_request_id" validate:"required,min=1"`
	PullRequestName string `json:"pull_request_name" validate:"required,min=1"`
	AuthorID        string `json:"author_id" validate:"required,min=1"`
}

type PRMergeRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required,min=1"`
}

type PRReassignRequest struct {
	PullRequestID string `json:"pull_request_id" validate:"required,min=1"`
	OldUserID     string `json:"old_user_id" validate:"required,min=1"`
}
