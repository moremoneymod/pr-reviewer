package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

var (
	ErrPRNotFound      = errors.New("PR not found")
	ErrTeamNotFound    = errors.New("team not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrPRMerged        = errors.New("PR merged")
	ErrPRExists        = errors.New("PR exists")
	ErrTeamExists      = errors.New("team already exists")
	ErrNoCandidates    = errors.New("no candidates")
	ErrUserNotReviewer = errors.New("user not reviewer")
)

type PRProvider interface {
	Create(ctx context.Context, pr domain.PR) (*domain.PR, error)
	Get(ctx context.Context, prId string) (*domain.PR, error)
	Merge(ctx context.Context, prId string) (*domain.PR, error)
	GetPullRequestsIdsByReviewer(ctx context.Context, reviewerId string) ([]string, error)
	GetAllPR(ctx context.Context) ([]*domain.PR, error)
	GetPRStatistics(ctx context.Context) (*domain.PRStatistics, error)
}

type TeamProvider interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	GetTeamById(ctx context.Context, teamId int) (*domain.Team, error)
	GetAllTeam(ctx context.Context) ([]*domain.Team, error)
	GetTeamStatistics(ctx context.Context) (*domain.TeamStatistics, error)
}

type UserProvider interface {
	SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, error)
	GetReview(ctx context.Context, prIds []string) ([]*domain.PRShort, error)
	GetUser(ctx context.Context, userId string) (*domain.User, error)
	GetReviewers(ctx context.Context, teamId int, excludeUserIds []string, limit int) ([]string, error)
	ReplaceReviewer(ctx context.Context, newReviewerId string, oldReviewerId string, prId string) error
	GetUserStatistics(ctx context.Context) (*domain.UserStatistics, error)
	GetUserAssignmentStatistics(ctx context.Context) ([]domain.UserAssignmentStat, error)
}

type Service struct {
	log          *slog.Logger
	PRRepository PRProvider
	TeamProvider TeamProvider
	UserProvider UserProvider
}

func New(log *slog.Logger, prProvider PRProvider, teamProvider TeamProvider, userProvider UserProvider) *Service {
	return &Service{
		log:          log,
		PRRepository: prProvider,
		TeamProvider: teamProvider,
		UserProvider: userProvider,
	}
}
