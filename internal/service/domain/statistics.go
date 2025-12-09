package domain

type Statistics struct {
	UserAssignments []UserAssignmentStat
	TotalPRs        int
	OpenPRs         int
	MergedPRs       int
	TotalTeams      int
	TotalUsers      int
	ActiveUsers     int
}

type UserAssignmentStat struct {
	UserID            string
	Username          string
	TeamName          string
	TotalAssignments  int
	OpenAssignments   int
	MergedAssignments int
}
