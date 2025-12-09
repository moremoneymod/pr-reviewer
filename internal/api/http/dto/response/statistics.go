package response

type StatisticsResponse struct {
	Statistics StatisticsData `json:"statistics"`
}

type StatisticsData struct {
	UserAssignments map[string]int      `json:"user_assignments"`
	PRAssignments   map[string]int      `json:"pr_assignments"`
	TeamStats       map[string]TeamStat `json:"team_stats"`
	TotalPRs        int                 `json:"total_prs"`
	OpenPRs         int                 `json:"open_prs"`
	MergedPRs       int                 `json:"merged_prs"`
}

type TeamStat struct {
	MemberCount   int `json:"member_count"`
	ActiveMembers int `json:"active_members"`
	PRCount       int `json:"pr_count"`
}
