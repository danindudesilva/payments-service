package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danindudesilva/payments-service/internal/payments/domain"
	memoryrepo "github.com/danindudesilva/payments-service/internal/payments/repository/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stripe "github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

func TestStripeWebhook_AlreadyProcessedReturnsOK(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)
	repo := memoryrepo.NewRepository()

	processedRepo := newFakeProcessedWebhookEventRepo()
	processedRepo.processed["stripe:evt_already_done"] = true

	handler := newWebhookTestHandlerWithProcessedRepo(repo, now, processedRepo)

	payload := fmt.Sprintf(`{
		"id":"evt_already_done",
		"object":"event",
		"type":"payment_intent.succeeded",
		"api_version":"%s",
		"data":{"object":{"id":"pi_123","object":"payment_intent"}}
	}`, stripe.APIVersion)

	signature := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload: []byte(payload),
		Secret:  testWebhookSecret,
	})

	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewBufferString(payload))
	req.Header.Set("Stripe-Signature", signature.Header)
	res := httptest.NewRecorder()

	handler.handleStripeWebhook(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), `"already_processed":true`)
	assert.Equal(t, 0, processedRepo.saveCalls)
}

func TestStripeWebhook_HasProcessedEventFailureReturnsServerError(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)
	repo := memoryrepo.NewRepository()

	processedRepo := newFakeProcessedWebhookEventRepo()
	processedRepo.hasErr = errors.New("repo unavailable")

	handler := newWebhookTestHandlerWithProcessedRepo(repo, now, processedRepo)

	payload := fmt.Sprintf(`{
		"id":"evt_check_fail",
		"object":"event",
		"type":"charge.updated",
		"api_version":"%s",
		"data":{"object":{"id":"ch_123","object":"charge"}}
	}`, stripe.APIVersion)

	signature := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload: []byte(payload),
		Secret:  testWebhookSecret,
	})

	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewBufferString(payload))
	req.Header.Set("Stripe-Signature", signature.Header)
	res := httptest.NewRecorder()

	handler.handleStripeWebhook(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Contains(t, res.Body.String(), "failed to process webhook event")
}

func TestStripeWebhook_SaveProcessedEventFailureReturnsServerError(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)
	repo := memoryrepo.NewRepository()

	processedRepo := newFakeProcessedWebhookEventRepo()
	processedRepo.saveErr = errors.New("insert failed")

	handler := newWebhookTestHandlerWithProcessedRepo(repo, now, processedRepo)

	payload := fmt.Sprintf(`{
		"id":"evt_save_fail",
		"object":"event",
		"type":"charge.updated",
		"api_version":"%s",
		"data":{"object":{"id":"ch_123","object":"charge"}}
	}`, stripe.APIVersion)

	signature := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload: []byte(payload),
		Secret:  testWebhookSecret,
	})

	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewBufferString(payload))
	req.Header.Set("Stripe-Signature", signature.Header)
	res := httptest.NewRecorder()

	handler.handleStripeWebhook(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Contains(t, res.Body.String(), "failed to process webhook event")
}

func TestStripeWebhook_HandledEventProcessesAndMarksProcessed(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)
	repo := memoryrepo.NewRepository()

	attempt, err := domain.NewPaymentAttempt(
		"attempt_123",
		"order_123",
		"idem_123",
		"https://example.com/return",
		domain.Money{Amount: 2500, Currency: "GBP"},
		now,
	)
	require.NoError(t, err)

	err = attempt.LinkProvider("stripe", "pi_123", "secret_123", now)
	require.NoError(t, err)

	err = repo.Save(context.Background(), attempt)
	require.NoError(t, err)

	processedRepo := newFakeProcessedWebhookEventRepo()
	handler := newWebhookTestHandlerWithProcessedRepo(repo, now.Add(time.Minute), processedRepo)

	payload := fmt.Sprintf(`{
		"id":"evt_mark_processed",
		"object":"event",
		"type":"payment_intent.succeeded",
		"api_version":"%s",
		"data":{"object":{"id":"pi_123","object":"payment_intent"}}
	}`, stripe.APIVersion)

	signature := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload: []byte(payload),
		Secret:  testWebhookSecret,
	})

	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewBufferString(payload))
	req.Header.Set("Stripe-Signature", signature.Header)
	res := httptest.NewRecorder()

	handler.handleStripeWebhook(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, 1, processedRepo.saveCalls)

	got, err := repo.GetByID(context.Background(), "attempt_123")
	require.NoError(t, err)
	assert.Equal(t, domain.PaymentStatusSucceeded, got.Status)
}
