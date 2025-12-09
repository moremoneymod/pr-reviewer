package entity

type Statistics struct {
	UserAssignments map[string]int
	PRAssignments   map[string]int
	TeamStats       map[string]TeamStat
	TotalPRs        int
	OpenPRs         int
	MergedPRs       int
}

type TeamStat struct {
	MemberCount   int
	ActiveMembers int
	PRCount       int
}
