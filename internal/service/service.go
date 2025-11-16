package service

import (
	"context"
	"errors"
	"log/slog"

	serv "github.com/moremoneymod/pr-reviewer/internal/service/entity"
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
	Create(ctx context.Context, pr serv.PR) (*serv.PR, error)
	Get(ctx context.Context, prId string) (*serv.PR, error)
	Merge(ctx context.Context, prId string) (*serv.PR, error)
	GetPullRequestsIdsByReviewer(ctx context.Context, reviewerId string) ([]string, error)
	GetAllPR(ctx context.Context) ([]*serv.PR, error)
}

type TeamProvider interface {
	CreateTeam(ctx context.Context, team *serv.Team) (*serv.Team, error)
	GetTeam(ctx context.Context, teamName string) (*serv.Team, error)
	GetTeamById(ctx context.Context, teamId int) (*serv.Team, error)
	GetAllTeam(ctx context.Context) ([]*serv.Team, error)
}

type UserProvider interface {
	SetIsActive(ctx context.Context, userId string, isActive bool) (*serv.User, error)
	GetReview(ctx context.Context, prIds []string) ([]*serv.PRShort, error)
	GetUser(ctx context.Context, userId string) (*serv.User, error)
	GetReviewers(ctx context.Context, teamId int, excludeUserIds []string, limit int) ([]string, error)
	ReplaceReviewer(ctx context.Context, newReviewerId string, oldReviewerId string, prId string) error
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
