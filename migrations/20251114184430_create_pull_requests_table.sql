-- +goose Up
-- +goose StatementBegin
CREATE TABLE pull_requests (
                               id VARCHAR(50) PRIMARY KEY,
                               name VARCHAR(500) NOT NULL,
                               author_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
                               status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               merged_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pull_requests;
-- +goose StatementEnd
