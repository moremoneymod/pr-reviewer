package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/repository"
	serv "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

func (s *Service) Create(ctx context.Context, team *serv.Team) (*serv.Team, error) {
	const op = "internal.service.team.Create"

	log := s.log.With(
		slog.String("op", op),
		slog.String("teamName", team.Name))

	log.Info("attempting to create new team")
	teamEntity, err := s.TeamProvider.CreateTeam(ctx, team)
	if errors.Is(err, repository.ErrTeamExists) {
		log.Warn("team already exists")
		return nil, fmt.Errorf("%s: %w", op, ErrTeamExists)
	}
	if err != nil {
		log.Error("failed to create team", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully created team")
	return teamEntity, nil
}

func (s *Service) Get(ctx context.Context, teamName string) (*serv.Team, error) {
	const op = "internal.service.team.Get"

	log := s.log.With(
		slog.String("op", op),
		slog.String("teamName", teamName))

	log.Info("attempting to get team")
	team, err := s.TeamProvider.GetTeam(ctx, teamName)
	if errors.Is(err, repository.ErrTeamNotFound) {
		log.Warn("team not found")
		return nil, fmt.Errorf("%s: %w", op, ErrTeamNotFound)
	}
	if err != nil {
		log.Error("failed to get team", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully got team")
	return team, nil
}
