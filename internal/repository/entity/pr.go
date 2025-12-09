package entity

import "time"

type PR struct {
	CreatedAt time.Time  `db:"created_at"`
	MergedAt  *time.Time `db:"merged_at"`
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	AuthorID  string     `db:"author_id"`
	Status    string     `db:"status"`
	Reviewers []string   `db:"-"`
}

type PRShort struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	AuthorID string `db:"author_id"`
	Status   string `db:"status"`
}
