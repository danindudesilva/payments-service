package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/danindudesilva/payments-service/internal/config"
)

func NewRouter(cfg config.Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteMethodNotAllowed(w)
			return
		}

		WriteJSON(w, http.StatusOK, map[string]any{
			"status": "ok",
			"env":    cfg.AppEnv,
		})
	})

	return chain(
		mux,
		requestID(),
		timeout(),
		recoverPanic(logger),
		requestLogger(logger),
	)
}
