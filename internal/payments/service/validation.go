package service

import "strings"

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}

func (in CreatePaymentAttemptInput) Validate() error {
	if strings.TrimSpace(in.OrderID) == "" {
		return ValidationError{Field: "order_id", Message: "must not be empty"}
	}

	if strings.TrimSpace(in.IdempotencyKey) == "" {
		return ValidationError{Field: "idempotency_key", Message: "must not be empty"}
	}

	if in.Amount <= 0 {
		return ValidationError{Field: "amount", Message: "must be greater than zero"}
	}

	if strings.TrimSpace(in.Currency) == "" {
		return ValidationError{Field: "currency", Message: "must not be empty"}
	}

	if strings.TrimSpace(in.ReturnURL) == "" {
		return ValidationError{Field: "return_url", Message: "must not be empty"}
	}

	return nil
}
