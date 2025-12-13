package reassign

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/converter"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/request"
	apiErrors "github.com/moremoneymod/pr-reviewer/internal/errors"
	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	"github.com/moremoneymod/pr-reviewer/internal/service"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/domain"
)

type PRAssigner interface {
	Reassign(ctx context.Context, prId string, oldUserId string) (*domain.PR, error)
}

func New(log *slog.Logger, prAssigner PRAssigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.pullrequest.reassign.New"

		log := log.With(
			slog.String("op", op))

		var req request.PRReassignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("error decoding body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeBadRequest, "error decoding body"))

			return
		}

		log = log.With(
			slog.String("prId", req.PullRequestID))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiErrors.ValidationError(validateErr))

			return
		}

		reassignedPR, err := prAssigner.Reassign(r.Context(), req.PullRequestID, req.OldUserID)
		if errors.Is(err, service.ErrPRNotFound) {
			log.Warn("PR not found")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeNotFound, "PR Not Found"))

			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			log.Warn("User not found")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeNotFound, "User Not Found"))
		}
		if errors.Is(err, service.ErrPRMerged) {
			log.Warn("PR merged")
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodePRMerged, "PR merged"))

			return
		}
		if errors.Is(err, service.ErrNoCandidates) {
			log.Warn("no candidates")
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeNoCandidate, "no candidates"))

			return
		}
		if errors.Is(err, service.ErrUserNotReviewer) {
			log.Warn("user not reviewer")
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeNotAssigned, "user not reviewer"))

			return
		}
		if err != nil {
			log.Error("error reassigning PR", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeInternalServer, "error reassigning PR"))
			return
		}

		response := converter.ToDTOPRFromDomain(reassignedPR)

		log.Info("pr reassigned successfully")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)

	}
}
