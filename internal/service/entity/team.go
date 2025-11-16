package entity

type Team struct {
	ID      int
	Name    string
	Members []Member
}

type Member struct {
	UserID   string
	Username string
	TeamID   int
	IsActive bool
}
