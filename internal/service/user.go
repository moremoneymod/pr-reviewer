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

func (s *Service) SetIsActive(ctx context.Context, userId string, isActive bool) (*serv.User, error) {
	const op = "internal.service.user.SetIsActive"

	log := s.log.With(
		slog.String("op", op),
		slog.String("userId", userId))

	log.Info("attempting to set user active flag")
	user, err := s.UserProvider.SetIsActive(ctx, userId, isActive)
	if errors.Is(err, repository.ErrUserNotFound) {
		log.Warn("user not found", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	if err != nil {
		log.Error("failed to set user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully set user active flag", slog.Any("isActive", isActive))
	return user, nil
}

func (s *Service) GetReview(ctx context.Context, userId string) ([]*serv.PRShort, error) {
	const op = "internal.service.getReview"

	log := s.log.With(
		slog.String("op", op),
		slog.String("userId", userId))

	prIds, err := s.PRRepository.GetPullRequestsIdsByReviewer(ctx, userId)
	if err != nil {
		log.Error("failed to get review", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("attempting to get review")
	prs, err := s.UserProvider.GetReview(ctx, prIds)
	if err != nil {
		log.Error("failed to get review", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully got review")
	return prs, nil
}
