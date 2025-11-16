package create

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/converter"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/request"
	errors2 "github.com/moremoneymod/pr-reviewer/internal/errors"
	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/service"
	serv "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

type PRCreator interface {
	CreatePR(ctx context.Context, prId string, prName string, authorId string) (*serv.PR, error)
}

func New(log *slog.Logger, prCreator PRCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.pullrequest.create.New"

		log := log.With(
			slog.String("op", op))

		var req request.PRCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("error decoding body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, errors2.NewErrorResponse(errors2.ErrorCodeBadRequest, "error decoding body"))

			return
		}

		log = log.With(
			slog.String("prId", req.PullRequestID))

		createdPR, err := prCreator.CreatePR(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
		if errors.Is(err, service.ErrPRExists) {
			log.Warn("PR already exists", sl.Err(err))
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, errors2.NewErrorResponse(errors2.ErrorCodePRExists, "PR already exists"))

			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			log.Warn("User not found", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, errors2.NewErrorResponse(errors2.ErrorCodeNotFound, "User not found"))

			return
		}
		if errors.Is(err, service.ErrTeamNotFound) {
			log.Warn("Team not found", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, errors2.NewErrorResponse(errors2.ErrorCodeNotFound, "Team not found"))

			return
		}
		if err != nil {
			log.Error("error creating PR", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, errors2.NewErrorResponse(errors2.ErrorCodeInternalServer, "error creating PR"))

			return
		}

		response := converter.ToPRDtoFromService(createdPR)

		log.Info("successfully created PR", slog.String("id", createdPR.ID))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)

		return
	}
}
