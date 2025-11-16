package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/pullrequest/create"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/pullrequest/merge"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/pullrequest/reassign"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/statistic"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/team/add"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/team/get"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/users/get_review"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/handlers/users/set_active"
	"github.com/moremoneymod/pr-reviewer/internal/config"
	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/service"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	service    *service.Service
}

func New(log *slog.Logger, httpConfig config.HTTPConfig, service *service.Service) *App {

	router := setupRouter(log, service)

	httpServer := &http.Server{
		Addr:         httpConfig.Address(),
		Handler:      router,
		ReadTimeout:  httpConfig.Timeout(),
		WriteTimeout: httpConfig.Timeout(),
		IdleTimeout:  httpConfig.IDLETimeout(),
	}

	return &App{
		log:        log,
		httpServer: httpServer,
		service:    service,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "internal.app.http.Run"

	log := app.log.With(
		slog.String("op", op))

	if err := app.httpServer.ListenAndServe(); err != nil {
		log.Error("failed to start http server", sl.Err(err))
		return err
	}
	log.Error("stopped http server")
	return nil
}

func setupRouter(log *slog.Logger, service *service.Service) *chi.Mux {
	const op = "internal.app.http.setupRouter"

	router := chi.NewRouter()
	router.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", create.New(log, service))
		r.Post("/merge", merge.New(log, service))
		r.Post("/reassign", reassign.New(log, service))
	})
	router.Route("/team", func(r chi.Router) {
		r.Post("/add", add.New(log, service))
		r.Get("/get", get.New(log, service))
	})
	router.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", set_active.New(log, service))
		r.Get("/getReview", get_review.New(log, service))
	})
	router.Get("/statistics", statistic.New(log, service))
	return router
}
