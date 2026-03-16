package postgres

import (
	"context"
	"testing"

	"github.com/danindudesilva/payments-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessedWebhookEventRepository_SaveAndHasProcessedEvent(t *testing.T) {
	pool := testutil.NewTestPool(t)
	repo := NewProcessedWebhookEventRepository(pool)

	processed, err := repo.HasProcessedEvent(context.Background(), "stripe", "evt_123")
	require.NoError(t, err)
	assert.False(t, processed)

	err = repo.SaveProcessedEvent(context.Background(), "stripe", "evt_123", "payment_intent.succeeded")
	require.NoError(t, err)

	processed, err = repo.HasProcessedEvent(context.Background(), "stripe", "evt_123")
	require.NoError(t, err)
	assert.True(t, processed)
}

func TestProcessedWebhookEventRepository_SaveProcessedEvent_DuplicateReturnsError(t *testing.T) {
	pool := testutil.NewTestPool(t)
	repo := NewProcessedWebhookEventRepository(pool)

	err := repo.SaveProcessedEvent(context.Background(), "stripe", "evt_duplicate", "payment_intent.succeeded")
	require.NoError(t, err)

	err = repo.SaveProcessedEvent(context.Background(), "stripe", "evt_duplicate", "payment_intent.succeeded")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "save processed webhook event")
}

func TestProcessedWebhookEventRepository_HasProcessedEvent_IsProviderScoped(t *testing.T) {
	pool := testutil.NewTestPool(t)
	repo := NewProcessedWebhookEventRepository(pool)

	err := repo.SaveProcessedEvent(context.Background(), "stripe", "evt_same", "payment_intent.succeeded")
	require.NoError(t, err)

	processed, err := repo.HasProcessedEvent(context.Background(), "stripe", "evt_same")
	require.NoError(t, err)
	assert.True(t, processed)

	processed, err = repo.HasProcessedEvent(context.Background(), "adyen", "evt_same")
	require.NoError(t, err)
	assert.False(t, processed)
}
