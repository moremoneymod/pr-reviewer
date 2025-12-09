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

func ToDomainTeamFromEntity(TeamEntity *entity.Team) *domain.Team {
	team := &domain.Team{
		ID:   TeamEntity.ID,
		Name: TeamEntity.Name,
	}

	team.Members = make([]domain.Member, len(TeamEntity.Members))

	for i, member := range TeamEntity.Members {
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

func ToDomainUserFromEntity(UserEntity *entity.User) *domain.User {
	return &domain.User{
		ID:       UserEntity.ID,
		Username: UserEntity.Username,
		TeamID:   UserEntity.TeamID,
		TeamName: UserEntity.TeamName,
		IsActive: UserEntity.IsActive,
	}
}

func ToDomainMembersFromEntity(MembersEntity []entity.Member) []domain.Member {
	members := make([]domain.Member, len(MembersEntity))
	for i, member := range MembersEntity {
		members[i] = ToDomainMemberFromEntity(member)
	}

	return members
}

func ToDomainUserAssignmentStatsFromEntity(UserAssignmentStatsEntity []*entity.UserAssignmentStatistics) []*domain.UserAssignmentStat {
	userAssignmentStats := make([]*domain.UserAssignmentStat, len(UserAssignmentStatsEntity))
	for i, userAssignmentStat := range UserAssignmentStatsEntity {
		userAssignmentStats[i] = &domain.UserAssignmentStat{
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

func ToDomainUserStatisticsFromEntity(UserStatisticsEntity *entity.UserStatistics) *domain.UserStatistics {
	return &domain.UserStatistics{
		TotalUsers:  UserStatisticsEntity.TotalUsers,
		ActiveUsers: UserStatisticsEntity.ActiveUsers,
	}
}

func ToDomainTeamStatisticsFromEntity(TeamStatisticsEntity *entity.TeamStatistics) *domain.TeamStatistics {
	return &domain.TeamStatistics{
		TotalTeams: TeamStatisticsEntity.TotalTeams,
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
