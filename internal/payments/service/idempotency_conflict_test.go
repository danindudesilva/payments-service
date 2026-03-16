package service

import (
	"context"
	"testing"
	"time"

	"github.com/danindudesilva/payments-service/internal/payments/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type conflictRepo struct {
	existing    *domain.PaymentAttempt
	saveCalls   int
	lookupCalls int
}

func (r *conflictRepo) Save(ctx context.Context, attempt *domain.PaymentAttempt) error {
	r.saveCalls++
	return &pgconn.PgError{Code: "23505"}
}

func (r *conflictRepo) GetByID(ctx context.Context, id string) (*domain.PaymentAttempt, error) {
	return nil, domain.ErrPaymentNotFound
}

func (r *conflictRepo) GetByProviderPaymentID(ctx context.Context, providerPaymentID string) (*domain.PaymentAttempt, error) {
	return nil, domain.ErrPaymentNotFound
}

func (r *conflictRepo) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.PaymentAttempt, error) {
	r.lookupCalls++

	if r.lookupCalls == 1 {
		return nil, domain.ErrPaymentNotFound
	}

	return r.existing, nil
}

func TestService_CreatePaymentAttempt_UniqueConflictReturnsExistingAttempt(t *testing.T) {
	now := time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC)

	existing, err := domain.NewPaymentAttempt(
		"attempt_existing",
		"order_123",
		"idem_123",
		"https://example.com/return",
		domain.Money{Amount: 2500, Currency: "GBP"},
		now,
	)
	require.NoError(t, err)

	repo := &conflictRepo{
		existing: existing,
	}

	gatewayCalls := 0
	gateway := &fakeGateway{
		createPaymentFunc: func(ctx context.Context, request domain.CreateProviderPaymentRequest) (domain.CreateProviderPaymentResult, error) {
			gatewayCalls++
			return domain.CreateProviderPaymentResult{
				ProviderName:      "stripe",
				ProviderPaymentID: "pi_123",
				ClientSecret:      "secret_123",
				Status:            domain.PaymentStatusPending,
			}, nil
		},
		getPaymentFunc: func(ctx context.Context, providerPaymentID string) (domain.CreateProviderPaymentResult, error) {
			return domain.CreateProviderPaymentResult{}, nil
		},
	}

	svc := New(
		repo,
		gateway,
		func() time.Time { return now },
		func() string { return "attempt_new" },
	)

	result, err := svc.CreatePaymentAttempt(context.Background(), CreatePaymentAttemptInput{
		OrderID:        "order_123",
		IdempotencyKey: "idem_123",
		Amount:         2500,
		Currency:       "GBP",
		ReturnURL:      "https://example.com/return",
		Description:    "test payment",
	})
	require.NoError(t, err)

	require.NotNil(t, result)
	assert.Equal(t, "attempt_existing", result.Attempt.ID)
	assert.Equal(t, 1, repo.saveCalls)
	assert.Equal(t, 2, repo.lookupCalls)
	assert.Equal(t, 1, gatewayCalls)
}
