package entity

type Statistics struct {
	TotalPRs        int
	OpenPRs         int
	MergedPRs       int
	UserAssignments map[string]int
	PRAssignments   map[string]int
	TeamStats       map[string]TeamStat
}

type TeamStat struct {
	MemberCount   int
	ActiveMembers int
	PRCount       int
}
