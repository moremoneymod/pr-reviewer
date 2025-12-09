package converter

import (
	"github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func ToDomainPRFromEntity(PREntity *entity.PR) *domain.PR {
	return &domain.PR{
		ID:        PREntity.ID,
		Name:      PREntity.Name,
		AuthorID:  PREntity.AuthorID,
		Status:    StringToPRStatus(PREntity.Status),
		Reviewers: PREntity.Reviewers,
		CreatedAt: &PREntity.CreatedAt,
		MergedAt:  PREntity.MergedAt,
	}
}

func ToDomainTeamFromEntity(teamEntity *entity.Team) *domain.Team {
	team := &domain.Team{
		ID:   teamEntity.ID,
		Name: teamEntity.Name,
	}

	team.Members = make([]domain.Member, len(teamEntity.Members))

	for i, member := range teamEntity.Members {
		team.Members[i] = ToDomainMemberFromEntity(member)
	}

	return team
}

func ToDomainMemberFromEntity(memberEntity entity.Member) domain.Member {
	return domain.Member{
		UserID:   memberEntity.UserID,
		Username: memberEntity.Username,
		TeamID:   memberEntity.TeamID,
		IsActive: memberEntity.IsActive,
	}
}

func ToDomainPRShortsFromEntity(PRsEntity []*entity.PRShort) []*domain.PRShort {
	prShorts := make([]*domain.PRShort, len(PRsEntity))
	for i, pr := range PRsEntity {
		prShorts[i] = &domain.PRShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		}
	}

	return prShorts
}

func ToDomainUserFromEntity(userEntity *entity.User) *domain.User {
	return &domain.User{
		ID:       userEntity.ID,
		Username: userEntity.Username,
		TeamID:   userEntity.TeamID,
		TeamName: userEntity.TeamName,
		IsActive: userEntity.IsActive,
	}
}

func ToDomainMembersFromEntity(membersEntity []entity.Member) []domain.Member {
	members := make([]domain.Member, len(membersEntity))
	for i, member := range membersEntity {
		members[i] = ToDomainMemberFromEntity(member)
	}

	return members
}

func ToDomainUserAssignmentStatsFromEntity(
	userAssignmentStatsEntity []entity.UserAssignmentStatistics,
) []domain.UserAssignmentStat {
	userAssignmentStats := make([]domain.UserAssignmentStat, len(userAssignmentStatsEntity))
	for i, userAssignmentStat := range userAssignmentStatsEntity {
		userAssignmentStats[i] = domain.UserAssignmentStat{
			UserID:            userAssignmentStat.UserID,
			Username:          userAssignmentStat.Username,
			TeamName:          userAssignmentStat.TeamName,
			TotalAssignments:  userAssignmentStat.TotalAssignments,
			OpenAssignments:   userAssignmentStat.OpenAssignments,
			MergedAssignments: userAssignmentStat.MergedAssignments,
		}
	}

	return userAssignmentStats
}

func ToDomainPRStatisticsFromEntity(PRStatisticsEntity *entity.PRStatistics) *domain.PRStatistics {
	return &domain.PRStatistics{
		TotalPRs:  PRStatisticsEntity.TotalPRs,
		OpenPRs:   PRStatisticsEntity.OpenPRs,
		MergedPRs: PRStatisticsEntity.MergedPRs,
	}
}

func ToDomainUserStatisticsFromEntity(
	userStatisticsEntity *entity.UserStatistics,
) *domain.UserStatistics {
	return &domain.UserStatistics{
		TotalUsers:  userStatisticsEntity.TotalUsers,
		ActiveUsers: userStatisticsEntity.ActiveUsers,
	}
}

func ToDomainTeamStatisticsFromEntity(
	teamStatisticsEntity *entity.TeamStatistics,
) *domain.TeamStatistics {
	return &domain.TeamStatistics{
		TotalTeams: teamStatisticsEntity.TotalTeams,
	}
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
