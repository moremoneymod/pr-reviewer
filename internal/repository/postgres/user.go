package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/moremoneymod/pr-reviewer/internal/repository"
	"github.com/moremoneymod/pr-reviewer/internal/repository/converter"
	entity "github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func (s *Storage) SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, error) {
	const op = "internal.repository.postgres.user.SetIsActive"

	builder := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("is_active", isActive).
		Where(sq.Eq{"id": userId})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() != 1 {
		return nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
	}
	user, err := s.GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) GetReview(ctx context.Context, prIds []string) ([]*domain.PRShort, error) {
	const op = "internal.repository.postgres.user.GetReview"

	builder := sq.Select("id", "name", "author_id", "status").
		PlaceholderFormat(sq.Dollar).
		From("pull_requests").
		Where(sq.Eq{"id": prIds})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s%w", op, err)
	}

	var result []*entity.PRShort
	err = pgxscan.Select(ctx, s.pgxPool, &result, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ToDomainPRShortsFromEntity(result), err
}

func (s *Storage) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	const op = "internal.repository.postgres.user.GetUser"

	builder := sq.Select("id", "username", "team_id", "is_active").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": userId})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var result entity.User
	err = pgxscan.Get(ctx, s.pgxPool, &result, query, args...)
	if pgxscan.NotFound(err) {
		return nil, repository.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ToDomainUserFromEntity(&result), nil
}

func (s *Storage) GetReviewers(ctx context.Context, teamId int, excludeUserIds []string, limit int) ([]string, error) {
	const op = "internal.repository.postgres.user.GetReviewers"

	builder := sq.Select("id").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"team_id": teamId}, sq.Eq{"is_active": true})

	if len(excludeUserIds) > 0 {
		for _, excludeUserId := range excludeUserIds {
			builder = builder.Where(sq.NotEq{"id": excludeUserId})
		}
	}

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}

	builder = builder.OrderBy("RANDOM()")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var result []string
	err = pgxscan.Select(ctx, s.pgxPool, &result, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *Storage) ReplaceReviewer(ctx context.Context, newReviewerId string, oldReviewerId string, prId string) error {
	const op = "internal.repository.postgres.user.ReplaceReviewer"

	builder := sq.Update("pr_reviewers").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"pr_id": prId}).
		Where(sq.Eq{"user_id": oldReviewerId}).
		Set("user_id", newReviewerId)
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = s.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserStatistics(ctx context.Context) (*domain.UserStatistics, error) {
	const op = "internal.repository.postgres.user.GetUserStatistics"

	builder := sq.Select("count(*) as total, count(CASE WHEN is_active = true THEN 1 END) as active").
		From("users")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var userStatistics entity.UserStatistics
	err = pgxscan.Get(ctx, s.pgxPool, &userStatistics, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ToDomainUserStatisticsFromEntity(&userStatistics), nil
}
