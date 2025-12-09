package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func (s *Service) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	const op = "internal.service.statistics.GetStatistics"

	log := s.log.With(
		slog.String("op", op))

	prStatistics, err := s.PRRepository.GetPRStatistics(ctx)
	if err != nil {
		log.Error("failed to get pr statistics", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userStatistics, err := s.UserProvider.GetUserStatistics(ctx)
	if err != nil {
		log.Error("failed to get user statistics", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	teamStatistics, err := s.TeamProvider.GetTeamStatistics(ctx)
	if err != nil {
		log.Error("failed to get team statistics", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userAssignmentsStat, err := s.UserProvider.GetUserAssignmentStatistics(ctx)
	if err != nil {
		log.Error("failed to get user assignment statistics", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	statistics := domain.Statistics{
		UserAssignments: userAssignmentsStat,
		TotalPRs:        prStatistics.TotalPRs,
		OpenPRs:         prStatistics.OpenPRs,
		MergedPRs:       prStatistics.MergedPRs,
		TotalTeams:      teamStatistics.TotalTeams,
		TotalUsers:      userStatistics.TotalUsers,
		ActiveUsers:     userStatistics.ActiveUsers,
	}

	return &statistics, nil
}
