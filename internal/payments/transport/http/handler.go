package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	basehttp "github.com/danindudesilva/payments-service/internal/httpserver"
	"github.com/danindudesilva/payments-service/internal/payments/domain"
	"github.com/danindudesilva/payments-service/internal/payments/service"
)

type Handler struct {
	service *service.Service
	logger  *slog.Logger
}

func NewHandler(service *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/payment-attempts", h.handlePaymentAttempts)
	mux.HandleFunc("/payment-attempts/", h.handlePaymentAttemptByID)
}

func (h *Handler) handlePaymentAttempts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createPaymentAttempt(w, r)
	default:
		basehttp.WriteMethodNotAllowed(w)
	}
}

func (h *Handler) handlePaymentAttemptByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getPaymentAttempt(w, r)
	default:
		basehttp.WriteMethodNotAllowed(w)
	}
}

func (h *Handler) createPaymentAttempt(w http.ResponseWriter, r *http.Request) {
	var request createPaymentAttemptRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		basehttp.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid json body",
		})
		return
	}

	result, err := h.service.CreatePaymentAttempt(r.Context(), service.CreatePaymentAttemptInput{
		OrderID:     request.OrderID,
		Amount:      request.Amount,
		Currency:    request.Currency,
		ReturnURL:   request.ReturnURL,
		Description: request.Description,
	})
	if err != nil {
		h.logger.Error("create payment attempt failed",
			slog.String("request_id", basehttp.RequestIDFromContext(r.Context())),
			slog.String("error", err.Error()),
		)

		basehttp.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	basehttp.WriteJSON(w, http.StatusCreated, toPaymentAttemptResponse(result.Attempt))
}

func (h *Handler) getPaymentAttempt(w http.ResponseWriter, r *http.Request) {
	attemptID := strings.TrimPrefix(r.URL.Path, "/payment-attempts/")
	if strings.TrimSpace(attemptID) == "" {
		basehttp.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "payment attempt id is required",
		})
		return
	}

	attempt, err := h.service.GetPaymentAttempt(r.Context(), attemptID)
	if err != nil {
		h.logger.Error("get payment attempt failed",
			slog.String("request_id", basehttp.RequestIDFromContext(r.Context())),
			slog.String("attempt_id", attemptID),
			slog.String("error", err.Error()),
		)

		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrPaymentNotFound) {
			status = http.StatusNotFound
		}

		basehttp.WriteJSON(w, status, map[string]any{
			"error": err.Error(),
		})
		return
	}

	basehttp.WriteJSON(w, http.StatusOK, toPaymentAttemptResponse(attempt))
}
