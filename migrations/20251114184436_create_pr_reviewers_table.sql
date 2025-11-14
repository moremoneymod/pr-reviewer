-- +goose Up
-- +goose StatementBegin
CREATE TABLE pr_reviewers (
                              pr_id VARCHAR(50) REFERENCES pull_requests(id) ON DELETE CASCADE,
                              user_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
                              assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              PRIMARY KEY (pr_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pr_reviewers;
-- +goose StatementEnd
