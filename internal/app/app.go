package app

import (
	"context"
	"log/slog"

	"github.com/moremoneymod/pr-reviewer/internal/app/http"
	"github.com/moremoneymod/pr-reviewer/internal/config"
	"github.com/moremoneymod/pr-reviewer/internal/repository/postgres"
	"github.com/moremoneymod/pr-reviewer/internal/service"
)

type App struct {
	HTTPSrv *http.App
}

func New(ctx context.Context, log *slog.Logger, pgConfig string, httpConfig config.HTTPConfig) *App {
	repository, err := postgres.New(ctx, pgConfig)
	if err != nil {
		panic(err)
	}

	appService := service.New(log, repository, repository, repository)
	httpApp := http.New(log, httpConfig, appService)

	return &App{
		HTTPSrv: httpApp,
	}
}
