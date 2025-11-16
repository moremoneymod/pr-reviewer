package entity

import "time"

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	TeamID    int       `db:"team_id"`
	TeamName  string    `db:"-"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}
