package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moremoneymod/pr-reviewer/internal/repository"
	"github.com/moremoneymod/pr-reviewer/internal/repository/converter"
	entity "github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func (s *Storage) Create(ctx context.Context, pr domain.PR) (*domain.PR, error) {
	const op = "internal.repository.postgres.postgres.Create"

	prEntity := entity.PR{
		ID:        pr.ID,
		Name:      pr.Name,
		AuthorID:  pr.AuthorID,
		Status:    converter.PRStatusToString(pr.Status),
		Reviewers: pr.Reviewers,
	}

	tx, err := s.pgxPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback(ctx)

	builder := sq.Insert("pull_requests").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "name", "author_id", "status").
		Values(prEntity.ID, prEntity.Name, prEntity.AuthorID, prEntity.Status).
		Suffix("RETURNING created_at")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&prEntity.CreatedAt)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrPRExists)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(pr.Reviewers) > 0 {
		reviewerBuilder := sq.Insert("pr_reviewers").
			PlaceholderFormat(sq.Dollar).
			Columns("pr_id", "user_id")

		for _, reviewer := range pr.Reviewers {
			reviewerBuilder = reviewerBuilder.Values(pr.ID, reviewer)
		}

		query, args, err := reviewerBuilder.ToSql()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ToDomainPRFromEntity(&prEntity), nil
}

func (s *Storage) Get(ctx context.Context, prId string) (*domain.PR, error) {
	const op = "internal.repository.postgres.postgres.Get"

	builder := sq.Select("id", "name", "author_id", "status", "created_at", "merged_at").
		PlaceholderFormat(sq.Dollar).
		From("pull_requests").
		Where(sq.Eq{"id": prId})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var pr entity.PR
	err = pgxscan.Get(ctx, s.pgxPool, &pr, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	reviewersBuilder := sq.Select("user_id").
		PlaceholderFormat(sq.Dollar).
		From("pr_reviewers").
		Where(sq.Eq{"pr_id": prId})
	query, args, err = reviewersBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var reviewers []string
	err = pgxscan.Select(ctx, s.pgxPool, &reviewers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pr.Reviewers = reviewers
	return converter.ToDomainPRFromEntity(&pr), nil
}

func (s *Storage) Merge(ctx context.Context, prId string) (*domain.PR, error) {
	const op = "internal.repository.postgres.postgres.Merge"

	builder := sq.Update("pull_requests").
		PlaceholderFormat(sq.Dollar).
		Set("status", "MERGED").
		Set("merged_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": prId})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("%s: %w", op, repository.ErrPRNotFound)
	}

	pr, err := s.Get(ctx, prId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pr, nil
}

func (s *Storage) GetPullRequestsIdsByReviewer(ctx context.Context, reviewerId string) ([]string, error) {
	const op = "internal.repository.postgres.postgres.GetPullRequestsIdsByReviewer"

	builder := sq.Select("pr_id").
		PlaceholderFormat(sq.Dollar).
		From("pr_reviewers").
		Where(sq.Eq{"user_id": reviewerId})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var pullRequestsIds []string
	err = pgxscan.Select(ctx, s.pgxPool, &pullRequestsIds, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pullRequestsIds, nil
}

func (s *Storage) GetPRStatistics(ctx context.Context) (*domain.PRStatistics, error) {
	const op = "internal.repository.postgres.postgres.GetPRStatistics"

	builder := sq.Select(
		"COUNT(*) as total",
		"COUNT(CASE WHEN status = 'OPEN' THEN 1 END) as open",
		"COUNT(CASE WHEN status = 'MERGED' THEN 1 END) as merged",
	).From("pull_requests")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var pullRequestsStats entity.PRStatistics
	err = pgxscan.Select(ctx, s.pgxPool, &pullRequestsStats, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrStatisticsNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ToDomainPRStatisticsFromEntity(&pullRequestsStats), nil
}
