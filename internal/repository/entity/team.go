package entity

import "time"

type Team struct {
	CreatedAt time.Time `db:"created_at"`
	Name      string    `db:"name"`
	Members   []Member  `db:"-"`
	ID        int       `db:"id"`
}

type Member struct {
	CreatedAt time.Time `db:"created_at"`
	UserID    string    `db:"id"`
	Username  string    `db:"username"`
	TeamID    int       `db:"team_id"`
	IsActive  bool      `db:"is_active"`
}
