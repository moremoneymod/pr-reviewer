package set_active

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
	domain "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

type UserActivityChanger interface {
	SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, error)
}

func New(log *slog.Logger, userActivitySetter UserActivityChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.users.get_review.New"

		log := log.With(
			slog.String("op", op))

		var req request.UserActiveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("error decoding body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeBadRequest, "error decoding body"))

			return
		}

		log = log.With(slog.String("userId", req.UserID))

		updatedUser, err := userActivitySetter.SetIsActive(r.Context(), req.UserID, req.IsActive)
		if errors.Is(err, service.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeNotFound, "user not found"))

			return
		}
		if err != nil {
			log.Error("error setting user activity", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeInternalServer, "error setting user activity"))
		}

		response := converter.ToDTOUserFromDomain(updatedUser)

		log.Info("user activity setting successfully")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)

		return
	}
}
