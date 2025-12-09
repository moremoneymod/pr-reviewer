package response

type StatisticsResponse struct {
	Statistics StatisticsData `json:"statistics"`
}

type StatisticsData struct {
	UserAssignments []UserAssignmentStat `json:"assignments"`
	TotalPRs        int                  `json:"totalPRs"`
	OpenPRs         int                  `json:"openPRs"`
	MergedPRs       int                  `json:"mergedPRs"`
	TotalTeams      int                  `json:"totalTeams"`
	TotalUsers      int                  `json:"totalUsers"`
	ActiveUsers     int                  `json:"activeUsers"`
}

type UserAssignmentStat struct {
	UserID            string `json:"userId"`
	Username          string `json:"username"`
	TeamName          string `json:"teamName"`
	TotalAssignments  int    `json:"totalAssignments"`
	OpenAssignments   int    `json:"openAssignments"`
	MergedAssignments int    `json:"mergedAssignments"`
}
