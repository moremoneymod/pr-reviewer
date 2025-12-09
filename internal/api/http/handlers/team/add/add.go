package add

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/converter"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/request"
	apiErrors "github.com/moremoneymod/pr-reviewer/internal/errors"
	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/service"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

type TeamAdder interface {
	Create(ctx context.Context, team *domain.Team) (*domain.Team, error)
}

func New(log *slog.Logger, teamAdder TeamAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.team.add.New"

		log := log.With(
			slog.String("op", op))

		var req request.TeamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("error decoding body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeBadRequest, "error decoding body"))

			return
		}

		log = log.With(slog.String("teamName", req.TeamName))

		team := converter.ToDomainTeamFromDTO(req)

		createdTeam, err := teamAdder.Create(context.Background(), team)
		if errors.Is(err, service.ErrTeamExists) {
			log.Info("team already exists")
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeTeamExists, fmt.Sprintf("%s already exists", team.Name)))

			return
		}
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeInternalServer, "internal server error"))

			return
		}

		response := converter.ToDTOTeamFromDomain(createdTeam)

		log.Info("team created")

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response)

		return
	}

}
