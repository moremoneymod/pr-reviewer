package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moremoneymod/pr-reviewer/internal/repository"
	"github.com/moremoneymod/pr-reviewer/internal/repository/converter"
	repo "github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	serv "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func (s *Storage) Create(ctx context.Context, pr serv.PR) (*serv.PR, error) {
	const op = "internal.repository.postgres.postgres.Create"

	prEntity := repo.PR{
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

func (s *Storage) Get(ctx context.Context, prId string) (*serv.PR, error) {
	const op = "internal.repository.postgres.postgres.Get"

	builder := sq.Select("id", "name", "author_id", "status", "created_at", "merged_at").
		PlaceholderFormat(sq.Dollar).
		From("pull_requests").
		Where(sq.Eq{"id": prId})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var pr repo.PR
	err = pgxscan.Get(ctx, s.pgxPool, &pr, query, args...)
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

func (s *Storage) GetAllPR(ctx context.Context) ([]*serv.PR, error) {
	const op = "internal.repository.postgres.postgres.GetAllPR"

	// Используем CTE для получения reviewers
	builder := sq.Select(
		"p.id",
		"p.name",
		"p.author_id",
		"p.status",
		"p.created_at",
		"p.merged_at",
		`(
				SELECT array_agg(user_id) 
				FROM pr_reviewers pr 
				WHERE pr.pr_id = p.id
			) as reviewers`,
	).
		PlaceholderFormat(sq.Dollar).
		From("pull_requests p").
		OrderBy("p.created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	type prRow struct {
		ID        string     `db:"id"`
		Name      string     `db:"name"`
		AuthorID  string     `db:"author_id"`
		Status    string     `db:"status"`
		CreatedAt time.Time  `db:"created_at"`
		MergedAt  *time.Time `db:"merged_at"`
		Reviewers []string   `db:"reviewers"`
	}

	var rows []prRow
	err = pgxscan.Select(ctx, s.pgxPool, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	prs := make([]*serv.PR, len(rows))
	for i, row := range rows {

		prs[i] = &serv.PR{
			ID:        row.ID,
			Name:      row.Name,
			AuthorID:  row.AuthorID,
			Status:    converter.StringToPRStatus(row.Status),
			Reviewers: row.Reviewers,
			CreatedAt: &row.CreatedAt,
			MergedAt:  row.MergedAt,
		}
	}

	return prs, nil
}

func (s *Storage) Merge(ctx context.Context, prId string) (*serv.PR, error) {
	const op = "internal.repository.postgres.postgres.Merge"

	builder := sq.Update("pull_requests").
		PlaceholderFormat(sq.Dollar).
		Set("status", "MERGED").
		Set("merged_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": prId})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s%w", op, err)
	}

	result, err := s.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s%w", op, err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("%s: %w", op, repository.ErrPRNotFound)
	}

	pr, err := s.Get(ctx, prId)
	if err != nil {
		return nil, fmt.Errorf("%s%w", op, err)
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
