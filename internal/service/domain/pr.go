package domain

import "time"

type PRStatus int

const (
	PRStatusOpen PRStatus = iota
	PRStatusMerged
)

type PR struct {
	CreatedAt *time.Time
	MergedAt  *time.Time
	ID        string
	Name      string
	AuthorID  string
	Reviewers []string
	Status    PRStatus
}
type PRShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}
