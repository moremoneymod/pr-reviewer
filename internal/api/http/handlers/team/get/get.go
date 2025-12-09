package get

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/converter"
	apiErrors "github.com/moremoneymod/pr-reviewer/internal/errors"
	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/service"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

type TeamService interface {
	Get(ctx context.Context, teamName string) (*domain.Team, error)
}

func New(log *slog.Logger, teamService TeamService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.New"

		log := log.With(
			slog.String("op", op))

		teamName := r.URL.Query().Get("team_name")

		if teamName == "" {
			log.Error("invalid request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeBadRequest, "team_name is required"))

			return
		}

		log = log.With(slog.String("teamName", teamName))

		team, err := teamService.Get(r.Context(), teamName)
		if errors.Is(err, service.ErrTeamNotFound) {
			log.Error("team not found")

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeNotFound, "team not found"))

			return
		}
		if err != nil {
			log.Error("internal error", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeInternalServer, "internal server error"))

			return
		}

		response := converter.ToDTOTeamFromDomain(team)

		log.Info("successfully read team", slog.String("team_name", teamName))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)

		return
	}
}
