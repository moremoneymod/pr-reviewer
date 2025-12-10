package health

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.handlers.health.New"

		log := log.With(slog.String("op", op))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"status": "ok"})
		log.Info("successfully health check")

		return
	}
}
