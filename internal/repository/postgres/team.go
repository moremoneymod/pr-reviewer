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
	"github.com/moremoneymod/pr-reviewer/internal/repository/entity"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

func (s *Storage) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	const op = "internal.repository.postgres.team.CreateTeam"

	tx, err := s.pgxPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	teamBuilder := sq.Insert("teams").
		PlaceholderFormat(sq.Dollar).
		Columns("name").
		Values(team.Name).
		Suffix("RETURNING id")

	teamQuery, args, err := teamBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.QueryRow(ctx, teamQuery, args...).Scan(&team.ID)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrTeamExists)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, member := range team.Members {
		userBuilder := sq.Insert("users").
			PlaceholderFormat(sq.Dollar).
			Columns("id", "username", "team_id", "is_active").
			Values(member.UserID, member.Username, team.ID, member.IsActive).
			Suffix(`
                ON CONFLICT (id) DO UPDATE SET
                    username = EXCLUDED.username,
                    team_id = EXCLUDED.team_id,
                    is_active = EXCLUDED.is_active;
            `)
		userQuery, args, err := userBuilder.ToSql()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		_, err = tx.Exec(ctx, userQuery, args...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return team, tx.Commit(ctx)
}

func (s *Storage) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	const op = "internal.repository.postgres.team.GetTeam"

	teamBuilder := sq.Select("id", "name", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("teams").
		Where(sq.Eq{"name": teamName})
	teamQuery, args, err := teamBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var team entity.Team
	err = pgxscan.Get(ctx, s.pgxPool, &team, teamQuery, args...)
	if pgxscan.NotFound(err) {
		return nil, repository.ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	membersBuilder := sq.Select("id", "username", "team_id", "is_active", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"team_id": team.ID})
	membersQuery, args, err := membersBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var members []entity.Member
	if err := pgxscan.Select(ctx, s.pgxPool, &members, membersQuery, args...); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	team.Members = members

	return converter.ToDomainTeamFromEntity(&team), nil
}

func (s *Storage) GetAllTeam(ctx context.Context) ([]*domain.Team, error) {
	const op = "internal.repository.postgres.team.GetAllTeam"

	builder := sq.Select(
		"t.id",
		"t.name",
		"t.created_at",
		"COALESCE(json_agg(json_build_object("+
			"'user_id', u.id, "+
			"'username', u.username, "+
			"'team_id', u.team_id, "+
			"'is_active', u.is_active, "+
			"'created_at', u.created_at"+
			")) FILTER (WHERE u.id IS NOT NULL), '[]') as members",
	).
		PlaceholderFormat(sq.Dollar).
		From("teams t").
		LeftJoin("users u ON t.id = u.team_id").
		GroupBy("t.id, t.name, t.created_at").
		OrderBy("t.created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	type teamRow struct {
		ID        int             `db:"id"`
		Name      string          `db:"name"`
		CreatedAt time.Time       `db:"created_at"`
		Members   []entity.Member `db:"members"`
	}

	var rows []teamRow
	err = pgxscan.Select(ctx, s.pgxPool, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	teams := make([]*domain.Team, len(rows))
	for i, row := range rows {
		teams[i] = &domain.Team{
			ID:      row.ID,
			Name:    row.Name,
			Members: converter.ToDomainMembersFromEntity(row.Members),
		}
	}

	return teams, nil
}

func (s *Storage) GetTeamById(ctx context.Context, teamId int) (*domain.Team, error) {
	const op = "internal.repository.postgres.team.GetTeam"

	teamBuilder := sq.Select("id", "name", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("teams").
		Where(sq.Eq{"id": teamId})
	teamQuery, args, err := teamBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var team entity.Team
	err = pgxscan.Get(ctx, s.pgxPool, &team, teamQuery, args...)
	if pgxscan.NotFound(err) {
		return nil, repository.ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	membersBuilder := sq.Select("id", "username", "team_id", "is_active", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"team_id": team.ID})
	membersQuery, args, err := membersBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var members []entity.Member
	if err := pgxscan.Select(ctx, s.pgxPool, &members, membersQuery, args...); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	team.Members = members

	return converter.ToDomainTeamFromEntity(&team), nil
}

func (s *Storage) GetTeamStatistics(ctx context.Context) (*domain.TeamStatistics, error) {
	const op = "internal.repository.postgres.team.GetTeamStatistics"

	builder := sq.Select("count(*) as total").From("teams")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var TeamStatistics entity.TeamStatistics
	err = pgxscan.Get(ctx, s.pgxPool, &TeamStatistics, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ToDomainTeamStatisticsFromEntity(&TeamStatistics), nil
}
