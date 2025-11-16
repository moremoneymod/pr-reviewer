package converter

import (
	"time"

	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/request"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/response"
	serv "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

func ToTeamFromDto(teamDto request.TeamRequest) *serv.Team {
	team := serv.Team{
		Name: teamDto.TeamName,
	}
	team.Members = make([]serv.Member, len(teamDto.Members))
	for i, member := range teamDto.Members {
		team.Members[i] = ToTeamMemberFromDto(member)
	}
	return &team
}

func ToTeamMemberFromDto(memberDto request.TeamMemberRequest) serv.Member {
	return serv.Member{
		UserID:   memberDto.UserID,
		Username: memberDto.Username,
		IsActive: memberDto.IsActive,
	}
}

func ToTeamDtoFromService(teamServ *serv.Team) response.TeamResponse {
	team := response.TeamResponse{
		TeamName: teamServ.Name,
	}
	team.Members = make([]response.TeamMember, 0, len(teamServ.Members))
	for _, member := range teamServ.Members {
		team.Members = append(team.Members, ToTeamMemberDtoFromService(member))
	}

	return team
}

func ToTeamMemberDtoFromService(teamMemberServ serv.Member) response.TeamMember {
	return response.TeamMember{
		UserID:   teamMemberServ.UserID,
		Username: teamMemberServ.Username,
		IsActive: teamMemberServ.IsActive,
	}
}

func ToPRDtoFromService(pr *serv.PR) response.PRResponse {
	var createdAtStr, mergedAtStr *string

	if pr.CreatedAt != nil {
		formatted := pr.CreatedAt.Format(time.RFC3339)
		createdAtStr = &formatted
	}

	if pr.MergedAt != nil {
		formatted := pr.MergedAt.Format(time.RFC3339)
		mergedAtStr = &formatted
	}

	prReviewers := pr.Reviewers
	if prReviewers == nil {
		prReviewers = []string{}
	}

	return response.PRResponse{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            PRStatusToString(pr.Status),
		AssignedReviewers: prReviewers,
		CreatedAt:         createdAtStr,
		MergedAt:          mergedAtStr,
	}
}

func ToUserDtoFromService(userServ *serv.User) response.UserResponse {
	return response.UserResponse{
		UserID:   userServ.ID,
		Username: userServ.Username,
		TeamName: userServ.TeamName,
		IsActive: userServ.IsActive,
	}
}

func ToPRsShortDtoFromService(prsShort []*serv.PRShort) []response.PRShortResponse {
	prsDto := make([]response.PRShortResponse, 0, len(prsShort))
	for _, prShort := range prsShort {
		prsDto = append(prsDto, ToPRShortDtoFromService(prShort))
	}
	return prsDto
}

func ToPRShortDtoFromService(prShort *serv.PRShort) response.PRShortResponse {
	return response.PRShortResponse{
		PullRequestID:   prShort.ID,
		PullRequestName: prShort.Name,
		AuthorID:        prShort.AuthorID,
		Status:          prShort.Status,
	}
}

func TOStatisticsDtoFromService(stats *serv.Statistics) response.StatisticsResponse {
	teamStats := make(map[string]response.TeamStat)
	for teamName, stat := range stats.TeamStats {
		teamStats[teamName] = response.TeamStat{
			MemberCount:   stat.MemberCount,
			ActiveMembers: stat.ActiveMembers,
			PRCount:       stat.PRCount,
		}
	}

	return response.StatisticsResponse{
		Statistics: response.StatisticsData{
			TotalPRs:        stats.TotalPRs,
			OpenPRs:         stats.OpenPRs,
			MergedPRs:       stats.MergedPRs,
			UserAssignments: stats.UserAssignments,
			PRAssignments:   stats.PRAssignments,
			TeamStats:       teamStats,
		},
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
