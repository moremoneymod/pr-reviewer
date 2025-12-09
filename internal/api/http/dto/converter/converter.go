package converter

import (
	"time"

	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/request"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/response"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func ToDomainTeamFromDTO(teamDTO request.TeamRequest) *domain.Team {
	members := make([]domain.Member, len(teamDTO.Members))
	for i, member := range teamDTO.Members {
		members[i] = ToDomainMemberFromDTO(member)
	}

	return &domain.Team{
		Name:    teamDTO.TeamName,
		Members: members,
	}
}

func ToDomainMemberFromDTO(memberDTO request.TeamMemberRequest) domain.Member {
	return domain.Member{
		UserID:   memberDTO.UserID,
		Username: memberDTO.Username,
		IsActive: memberDTO.IsActive,
	}
}

func ToDTOTeamFromDomain(teamDomain *domain.Team) response.TeamResponse {
	team := response.TeamResponse{
		TeamName: teamDomain.Name,
	}

	team.Members = make([]response.TeamMember, 0, len(teamDomain.Members))

	for _, member := range teamDomain.Members {
		team.Members = append(team.Members, ToDTOTeamMemberFromDomain(member))
	}

	return team
}

func ToDTOTeamMemberFromDomain(teamMemberDomain domain.Member) response.TeamMember {
	return response.TeamMember{
		UserID:   teamMemberDomain.UserID,
		Username: teamMemberDomain.Username,
		IsActive: teamMemberDomain.IsActive,
	}
}

func ToDTOPRFromDomain(PRDomain *domain.PR) response.PRResponse {
	var createdAtStr, mergedAtStr *string

	if PRDomain.CreatedAt != nil {
		formatted := PRDomain.CreatedAt.Format(time.RFC3339)
		createdAtStr = &formatted
	}

	if PRDomain.MergedAt != nil {
		formatted := PRDomain.MergedAt.Format(time.RFC3339)
		mergedAtStr = &formatted
	}

	prReviewers := PRDomain.Reviewers
	if prReviewers == nil {
		prReviewers = []string{}
	}

	return response.PRResponse{
		PullRequestID:     PRDomain.ID,
		PullRequestName:   PRDomain.Name,
		AuthorID:          PRDomain.AuthorID,
		Status:            PRStatusToString(PRDomain.Status),
		AssignedReviewers: prReviewers,
		CreatedAt:         createdAtStr,
		MergedAt:          mergedAtStr,
	}
}

func ToDTOUserFromDomain(userDomain *domain.User) response.UserResponse {
	return response.UserResponse{
		UserID:   userDomain.ID,
		Username: userDomain.Username,
		TeamName: userDomain.TeamName,
		IsActive: userDomain.IsActive,
	}
}

func ToDTOPRsShortFromDomain(PRsShortDomain []*domain.PRShort) []response.PRShortResponse {
	prsDto := make([]response.PRShortResponse, 0, len(PRsShortDomain))
	for _, prShort := range PRsShortDomain {
		prsDto = append(prsDto, ToDTOPRShortFromDomain(prShort))
	}

	return prsDto
}

func ToDTOPRShortFromDomain(PRShortDomain *domain.PRShort) response.PRShortResponse {
	return response.PRShortResponse{
		PullRequestID:   PRShortDomain.ID,
		PullRequestName: PRShortDomain.Name,
		AuthorID:        PRShortDomain.AuthorID,
		Status:          PRShortDomain.Status,
	}
}

func ToDTOStatisticsFromDomain(statisticsDomain *domain.Statistics) response.StatisticsResponse {
	teamStats := make(map[string]response.TeamStat)
	for teamName, stat := range statisticsDomain.TeamStats {
		teamStats[teamName] = response.TeamStat{
			MemberCount:   stat.MemberCount,
			ActiveMembers: stat.ActiveMembers,
			PRCount:       stat.PRCount,
		}
	}

	return response.StatisticsResponse{
		Statistics: response.StatisticsData{
			TotalPRs:        statisticsDomain.TotalPRs,
			OpenPRs:         statisticsDomain.OpenPRs,
			MergedPRs:       statisticsDomain.MergedPRs,
			UserAssignments: statisticsDomain.UserAssignments,
			PRAssignments:   statisticsDomain.PRAssignments,
			TeamStats:       teamStats,
		},
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
