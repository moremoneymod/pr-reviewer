package converter

import (
	"github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

func ToDomainPRFromEntity(pr *entity.PR) *domain.PR {
	return &domain.PR{
		ID:        pr.ID,
		Name:      pr.Name,
		AuthorID:  pr.AuthorID,
		Status:    StringToPRStatus(pr.Status),
		Reviewers: pr.Reviewers,
		CreatedAt: &pr.CreatedAt,
		MergedAt:  pr.MergedAt,
	}
}

func ToDomainTeamFromEntity(repoTeam *entity.Team) *domain.Team {
	team := &domain.Team{
		ID:   repoTeam.ID,
		Name: repoTeam.Name,
	}

	team.Members = make([]domain.Member, len(repoTeam.Members))

	for i, member := range repoTeam.Members {
		team.Members[i] = ToDomainMemberFromEntity(member)
	}

	return team
}

func ToDomainMemberFromEntity(member entity.Member) domain.Member {
	return domain.Member{
		UserID:   member.UserID,
		Username: member.Username,
		TeamID:   member.TeamID,
		IsActive: member.IsActive,
	}
}

func ToDomainPRShortsFromEntity(prs []*entity.PRShort) []*domain.PRShort {
	prShorts := make([]*domain.PRShort, len(prs))
	for i, pr := range prs {
		prShorts[i] = &domain.PRShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		}
	}

	return prShorts
}

func ToDomainUserFromEntity(repoUser *entity.User) *domain.User {
	return &domain.User{
		ID:       repoUser.ID,
		Username: repoUser.Username,
		TeamID:   repoUser.TeamID,
		TeamName: repoUser.TeamName,
		IsActive: repoUser.IsActive,
	}
}

func ToDomainMembersFromEntity(repoMembers []entity.Member) []domain.Member {
	members := make([]domain.Member, len(repoMembers))
	for i, member := range repoMembers {
		members[i] = ToDomainMemberFromEntity(member)
	}

	return members
}

func StringToPRStatus(status string) domain.PRStatus {
	switch status {
	case "OPEN":
		return domain.PRStatusOpen
	case "MERGED":
		return domain.PRStatusMerged
	default:
		return domain.PRStatusOpen
	}
}

func PRStatusToString(status domain.PRStatus) string {
	switch status {
	case domain.PRStatusOpen:
		return "OPEN"
	case domain.PRStatusMerged:
		return "MERGED"
	default:
		return "OPEN"
	}
}
