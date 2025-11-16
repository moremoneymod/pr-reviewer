package entity

import "time"

type PRStatus int

const (
	PRStatusOpen PRStatus = iota
	PRStatusMerged
)

type PR struct {
	ID        string
	Name      string
	AuthorID  string
	Status    PRStatus
	Reviewers []string
	CreatedAt *time.Time
	MergedAt  *time.Time
}
type PRShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}
