package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPool_InvalidURL(t *testing.T) {
	t.Parallel()

	pool, err := NewPool(context.Background(), Config{
		DatabaseURL: "://not-a-valid-url",
	})

	require.Error(t, err)
	require.Nil(t, pool)
}
