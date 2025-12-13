package repository

import "errors"

var (
	ErrTeamExists         = errors.New("team already exists")
	ErrTeamNotFound       = errors.New("team not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrPRExists           = errors.New("PR already exists")
	ErrPRNotFound         = errors.New("PR not found")
	ErrStatisticsNotFound = errors.New("statistics not found")
)
