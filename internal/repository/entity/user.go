package entity

import "time"

type User struct {
	CreatedAt time.Time `db:"created_at"`
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	TeamName  string    `db:"-"`
	TeamID    int       `db:"team_id"`
	IsActive  bool      `db:"is_active"`
}
