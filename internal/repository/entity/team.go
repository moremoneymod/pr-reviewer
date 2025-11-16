package entity

import "time"

type Team struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Members   []Member  `db:"-"`
	CreatedAt time.Time `db:"created_at"`
}

type Member struct {
	UserID    string    `db:"id"`
	Username  string    `db:"username"`
	TeamID    int       `db:"team_id"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}
