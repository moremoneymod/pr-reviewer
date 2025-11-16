package converter

import (
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/request"
	request2 "github.com/moremoneymod/pr-reviewer/internal/api/http/dto/response"
	repo "github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	serv "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

func ToTeamFromService(pr *serv.Team) *repo.Team {
	team := repo.Team{
		ID:   pr.ID,
		Name: pr.Name,
	}
	team.Members = make([]repo.Member, len(pr.Members))
	for i, member := range pr.Members {
		team.Members[i] = repo.Member{
			UserID:   member.UserID,
			Username: member.Username,
			TeamID:   member.TeamID,
			IsActive: member.IsActive,
		}
	}
	return &team
}

func ToUserDtoFromService(user *serv.User) request2.UserResponse {
	return request2.UserResponse{
		UserID:   user.ID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func ToTeamDtoFromService(team *serv.Team) request2.TeamResponse {
	teamDto := request2.TeamResponse{
		TeamName: team.Name,
	}
	teamDto.Members = make([]request2.TeamMember, 0, len(team.Members))
	for _, member := range team.Members {
		teamDto.Members = append(teamDto.Members, request2.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return teamDto
}

func ToMemberFromDto(memberDTO request.TeamMemberRequest, teamId string) serv.Member {
	return serv.Member{
		UserID:   memberDTO.UserID,
		Username: memberDTO.Username,
		TeamID:   teamId,
		IsActive: memberDTO.IsActive,
	}
}

func stringToPRStatus(status string) serv.PRStatus {
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
