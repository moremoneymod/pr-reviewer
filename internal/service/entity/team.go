package entity

type Team struct {
	Name    string
	Members []Member
	ID      int
}

type Member struct {
	UserID   string
	Username string
	TeamID   int
	IsActive bool
}
