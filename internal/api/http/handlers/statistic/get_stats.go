package statistic

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

type StatsProvider interface {
	GetStatistics(ctx context.Context) (*domain.Statistics, error)
}

func New(log *slog.Logger, statsProvider StatsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.statistics.get_stats.New"

		log := log.With(
			slog.String("op", op))

		stats, err := statsProvider.GetStatistics(r.Context())
		if err != nil {
			log.Error("failed getting stats", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, apiErrors.NewErrorResponse(apiErrors.ErrorCodeInternalServer, "failed to get statistics"))
			return
		}

		response := converter.ToDTOStatisticsFromDomain(stats)

		log.Info("success getting stats")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
		return
	}
}
