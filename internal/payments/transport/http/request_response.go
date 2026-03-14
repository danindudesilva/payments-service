package http

import "github.com/danindudesilva/payments-service/internal/payments/domain"

type createPaymentAttemptRequest struct {
	OrderID     string `json:"order_id"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	ReturnURL   string `json:"return_url"`
	Description string `json:"description"`
}

type paymentAttemptResponse struct {
	ID            string                `json:"id"`
	OrderID       string                `json:"order_id"`
	Status        string                `json:"status"`
	Amount        int64                 `json:"amount"`
	Currency      string                `json:"currency"`
	FailureReason string                `json:"failure_reason,omitempty"`
	NextAction    nextActionResponse    `json:"next_action"`
	Provider      providerResponse      `json:"provider"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
	CompletedAt   *string               `json:"completed_at,omitempty"`
}

type nextActionResponse struct {
	Type        string `json:"type"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

type providerResponse struct {
	Name            string `json:"name,omitempty"`
	PaymentID       string `json:"payment_id,omitempty"`
	ClientSecret    string `json:"client_secret,omitempty"`
}

func toPaymentAttemptResponse(attempt *domain.PaymentAttempt) paymentAttemptResponse {
	response := paymentAttemptResponse{
		ID:            attempt.ID,
		OrderID:       attempt.OrderID,
		Status:        string(attempt.Status),
		Amount:        attempt.Money.Amount,
		Currency:      attempt.Money.Currency,
		FailureReason: attempt.FailureReason,
		NextAction: nextActionResponse{
			Type:        string(attempt.NextAction.Type),
			RedirectURL: attempt.NextAction.RedirectURL,
		},
		Provider: providerResponse{
			Name:         attempt.Provider.ProviderName,
			PaymentID:    attempt.Provider.ProviderPaymentID,
			ClientSecret: attempt.Provider.ClientSecret,
		},
		CreatedAt: attempt.Timestamps.CreatedAt.UTC().Format(timeFormatRFC3339),
		UpdatedAt: attempt.Timestamps.UpdatedAt.UTC().Format(timeFormatRFC3339),
	}

	if attempt.Timestamps.CompletedAt != nil {
		completedAt := attempt.Timestamps.CompletedAt.UTC().Format(timeFormatRFC3339)
		response.CompletedAt = &completedAt
	}

	return response
}

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"
