package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	serv "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

func (s *Service) GetStatistics(ctx context.Context) (*serv.Statistics, error) {
	const op = "internal.service.statistics.GetStatistics"

	log := s.log.With(
		slog.String("op", op))

	allPRs, err := s.PRRepository.GetAllPR(ctx)
	if err != nil {
		log.Error("failed to get all PRs", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	allTeams, err := s.TeamProvider.GetAllTeam(ctx)
	if err != nil {
		log.Error("failed to get all teams", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	statistics := &serv.Statistics{
		UserAssignments: make(map[string]int),
		PRAssignments:   make(map[string]int),
		TeamStats:       make(map[string]serv.TeamStat),
	}

	for _, pr := range allPRs {
		statistics.TotalPRs++

		switch pr.Status {
		case serv.PRStatusOpen:
			statistics.OpenPRs++
		case serv.PRStatusMerged:
			statistics.MergedPRs++
		}

		statistics.PRAssignments[pr.ID] = len(pr.Reviewers)

		for _, reviewer := range pr.Reviewers {
			statistics.UserAssignments[reviewer]++
		}
	}

	for _, team := range allTeams {
		teamStat := serv.TeamStat{
			MemberCount:   len(team.Members),
			ActiveMembers: 0,
			PRCount:       0,
		}

		for _, member := range team.Members {
			if member.IsActive {
				teamStat.ActiveMembers++
			}
		}

		for _, pr := range allPRs {
			author, err := s.UserProvider.GetUser(ctx, pr.ID)
			if err == nil {
				if author.TeamName == team.Name {
					teamStat.PRCount++
				}
			}
		}
		statistics.TeamStats[team.Name] = teamStat
	}
	return statistics, nil
}
