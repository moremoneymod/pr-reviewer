package converter

import (
	repo "github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	serv "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

func ToPRFromRepository(pr *repo.PR) *serv.PR {
	return &serv.PR{
		ID:        pr.ID,
		Name:      pr.Name,
		AuthorID:  pr.AuthorID,
		Status:    StringToPRStatus(pr.Status),
		Reviewers: pr.Reviewers,
		CreatedAt: &pr.CreatedAt,
		MergedAt:  pr.MergedAt,
	}
}

func ToTeamFromRepository(repoTeam *repo.Team) *serv.Team {
	team := &serv.Team{
		ID:   repoTeam.ID,
		Name: repoTeam.Name,
	}
	team.Members = make([]serv.Member, len(repoTeam.Members))
	for i, member := range repoTeam.Members {
		team.Members[i] = ToMemberFromRepository(member)
	}
	return team
}

func ToMemberFromRepository(member repo.Member) serv.Member {
	return serv.Member{
		UserID:   member.UserID,
		Username: member.Username,
		TeamID:   member.TeamID,
		IsActive: member.IsActive,
	}
}

func ToPRShortsFromRepository(prs []*repo.PRShort) []*serv.PRShort {
	prShorts := make([]*serv.PRShort, len(prs))
	for i, pr := range prs {
		prShorts[i] = &serv.PRShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		}
	}
	return prShorts
}

func ToUserFromRepository(repoUser *repo.User) *serv.User {
	return &serv.User{
		ID:       repoUser.ID,
		Username: repoUser.Username,
		TeamID:   repoUser.TeamID,
		TeamName: repoUser.TeamName,
		IsActive: repoUser.IsActive,
	}
}

func ToMembersFromRepository(repoMembers []repo.Member) []serv.Member {
	members := make([]serv.Member, len(repoMembers))
	for i, member := range repoMembers {
		members[i] = ToMemberFromRepository(member)
	}
	return members
}

func StringToPRStatus(status string) serv.PRStatus {
	switch status {
	case "OPEN":
		return serv.PRStatusOpen
	case "MERGED":
		return serv.PRStatusMerged
	default:
		return serv.PRStatusOpen
	}
}

func PRStatusToString(status serv.PRStatus) string {
	switch status {
	case serv.PRStatusOpen:
		return "OPEN"
	case serv.PRStatusMerged:
		return "MERGED"
	default:
		return "OPEN"
	}
}
