package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/repository"
	"github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func (s *Service) CreatePR(ctx context.Context, prId string, prName string, authorId string) (*domain.PR, error) {
	const op = "internal.service.pr.CreatePR"

	log := s.log.With(
		slog.String("op", op),
		slog.String("prId", prId))

	log.Info("attempting to get user")
	author, err := s.UserProvider.GetUser(ctx, authorId)
	if errors.Is(err, repository.ErrUserNotFound) {
		log.Info("author not found")
		return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	if err != nil {
		log.Error("failed to get user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("attempting to get team")
	_, err = s.TeamProvider.GetTeamById(ctx, author.TeamID)
	if errors.Is(err, repository.ErrTeamNotFound) {
		log.Warn("team not found")
		return nil, fmt.Errorf("%s: %w", op, ErrTeamNotFound)
	}
	if err != nil {
		log.Error("failed to get team", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("attempting to get reviewers")
	reviewers, err := s.UserProvider.GetReviewers(ctx, author.TeamID, []string{authorId}, 2)
	if err != nil {
		log.Error("failed to get reviewers", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pr := domain.PR{
		ID:        prId,
		Name:      prName,
		AuthorID:  authorId,
		Status:    domain.PRStatusOpen,
		Reviewers: reviewers,
	}

	log.Info("attempting to create pr")
	prEntity, err := s.PRRepository.Create(ctx, pr)
	if errors.Is(err, repository.ErrPRExists) {
		log.Warn("pr already exists")
		return nil, fmt.Errorf("%s: %w", op, ErrPRExists)
	}
	if err != nil {
		log.Error("failed to create pr", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully created pr")
	return prEntity, nil
}

func (s *Service) Merge(ctx context.Context, prId string) (*domain.PR, error) {
	const op = "internal.service.pr.Merge"

	log := s.log.With(
		slog.String("op", op),
		slog.String("prId", prId))

	log.Info("attempting to merge pr")
	pr, err := s.PRRepository.Merge(ctx, prId)
	if errors.Is(err, repository.ErrPRNotFound) {
		log.Warn("pr not found")
		return nil, fmt.Errorf("%s: %w", op, ErrPRNotFound)
	}
	if err != nil {
		log.Error("failed to merge pr", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully merged pr")
	return pr, nil

}

func (s *Service) Reassign(ctx context.Context, prId string, oldUserId string) (*domain.PR, error) {
	const op = "internal.service.pr.Reassign"

	log := s.log.With(
		slog.String("op", op),
		slog.String("prId", prId),
		slog.String("oldUserId", oldUserId))

	log.Info("attempting to reassign pr")
	log.Info("attempting to get pr")
	pr, err := s.PRRepository.Get(ctx, prId)
	if errors.Is(err, repository.ErrPRNotFound) {
		log.Warn("pr not found")
		return nil, fmt.Errorf("%s: %w", op, ErrPRNotFound)
	}
	if err != nil {
		log.Error("failed to get pr", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if pr.Status == domain.PRStatusMerged {
		log.Warn("pr is already merged")
		return nil, fmt.Errorf("%s: %w", op, ErrPRMerged)
	}

	log.Info("attempting to get user")
	oldUser, err := s.UserProvider.GetUser(ctx, oldUserId)
	if errors.Is(err, repository.ErrUserNotFound) {
		log.Warn("user not found")
		return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	if err != nil {
		log.Error("failed to get user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !slices.Contains(pr.Reviewers, oldUserId) {
		log.Warn("user is not reviewer")
		return nil, fmt.Errorf("%s: %w", op, ErrUserNotReviewer)
	}

	log.Info("attempting to get reviewer candidates")
	excludeIds := pr.Reviewers
	excludeIds = append(excludeIds, pr.AuthorID)
	reviewerCandidates, err := s.UserProvider.GetReviewers(ctx, oldUser.TeamID, excludeIds, 1)
	if err != nil {
		log.Error("failed to get reviewer candidates", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(reviewerCandidates) == 0 {
		log.Warn("reviewer candidates not found")
		return nil, fmt.Errorf("%s: %w", op, ErrNoCandidates)
	}

	err = s.UserProvider.ReplaceReviewer(ctx, reviewerCandidates[0], oldUserId, prId)
	if err != nil {
		log.Error("failed to replace reviewer", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("attempting to get pr")
	newPr, err := s.PRRepository.Get(ctx, pr.ID)
	if err != nil {
		log.Error("failed to get pr", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully pr reassign")
	return newPr, nil
}
