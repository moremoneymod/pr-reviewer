-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       id VARCHAR(50) PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL,
                       is_active BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE USERS;
-- +goose StatementEnd
