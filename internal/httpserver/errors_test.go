package httpserver

import (
	"net/http"
	"testing"

	"github.com/danindudesilva/payments-service/internal/payments/domain"
	"github.com/danindudesilva/payments-service/internal/payments/service"
	"github.com/stretchr/testify/assert"
)

func TestMapError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "service validation error",
			err:        service.ValidationError{Field: "order_id", Message: "must not be empty"},
			wantStatus: http.StatusBadRequest,
			wantMsg:    "order_id: must not be empty",
		},
		{
			name:       "payment not found",
			err:        domain.ErrPaymentNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    "payment attempt not found",
		},
		{
			name:       "invalid transition",
			err:        domain.ErrInvalidTransition,
			wantStatus: http.StatusConflict,
			wantMsg:    "invalid payment state transition",
		},
		{
			name:       "provider already linked",
			err:        domain.ErrProviderAlreadyLinked,
			wantStatus: http.StatusConflict,
			wantMsg:    "provider payment is already linked",
		},
		{
			name:       "invalid money",
			err:        domain.ErrInvalidMoney,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "invalid payment request",
		},
		{
			name:       "invalid next action",
			err:        domain.ErrInvalidNextAction,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "invalid payment request",
		},
		{
			name:       "unknown error",
			err:        assert.AnError,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "internal server error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MapError(tt.err)

			assert.Equal(t, tt.wantStatus, got.StatusCode)
			assert.Equal(t, tt.wantMsg, got.Message)
		})
	}
}
