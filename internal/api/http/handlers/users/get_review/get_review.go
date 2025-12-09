package get_review

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/moremoneymod/pr-reviewer/internal/api/http/dto/converter"
	apiErrors "github.com/moremoneymod/pr-reviewer/internal/errors"
	"github.com/moremoneymod/pr-reviewer/internal/lib/logger/sl"
	domain "github.com/moremoneymod/pr-reviewer/internal/service/entity"
)

type UserReviewProvider interface {
	GetReview(ctx context.Context, userId string) ([]*domain.PRShort, error)
}

func New(log *slog.Logger, userReviewProvider UserReviewProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.user.get_review.New"

		log := log.With(
			slog.String("op", op))

		userId := r.URL.Query().Get("UserIdQuery")
		if userId == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeBadRequest, "missing userId"))
			return
		}

		log.With(slog.String("userId", userId))

		reviews, err := userReviewProvider.GetReview(r.Context(), userId)
		if err != nil {
			log.Error("failed to get reviews", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeInternalServer, "failed to get reviews"))

			return
		}

		response := converter.ToDTOPRsShortFromDomain(reviews)

		log.Info("successfully got reviews")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)

		return
	}
}
