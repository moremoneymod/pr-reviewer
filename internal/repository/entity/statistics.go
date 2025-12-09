package entity

type UserAssignmentStatistics struct {
	UserID            string `db:"user_id"`
	Username          string `db:"username"`
	TeamName          string `db:"team_name"`
	TotalAssignments  int    `db:"total_assignments"`
	OpenAssignments   int    `db:"open_assignments"`
	MergedAssignments int    `db:"merged_assignments"`
}

type PRStatistics struct {
	TotalPRs  int `db:"total"`
	OpenPRs   int `db:"open"`
	MergedPRs int `db:"merged"`
}

type UserStatistics struct {
	TotalUsers  int `db:"total"`
	ActiveUsers int `db:"active"`
}

type TeamStatistics struct {
	TotalTeams int `db:"total"`
}
