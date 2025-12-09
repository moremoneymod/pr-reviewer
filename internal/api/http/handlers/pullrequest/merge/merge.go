package merge

import (
	"context"
	"encoding/json"
	"errors"
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

type PRMerger interface {
	Merge(ctx context.Context, prId string) (*domain.PR, error)
}

func New(log *slog.Logger, prMerger PRMerger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.pullrequest.merge.New"

		log := log.With(
			slog.String("op", op))

		var req request.PRMergeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("error decoding body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeBadRequest, "error decoding body"))

			return
		}

		log = log.With(
			slog.String("prId", req.PullRequestID))

		mergedPR, err := prMerger.Merge(r.Context(), req.PullRequestID)
		if errors.Is(err, service.ErrPRNotFound) {
			log.Warn("PR not found", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeNotFound, "PR not found"))

			return
		}
		if err != nil {
			log.Error("error calling PRMerger", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeInternalServer, "error merging PR"))

			return
		}

		response := converter.ToDTOPRFromDomain(mergedPR)

		log.Info("pr merged successfully", slog.String("pr", mergedPR.ID))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
		return
	}
}
