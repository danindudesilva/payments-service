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
	require.ErrorContains(t, err, "parse pgx pool config")
}

func TestNewPool_PingFails(t *testing.T) {
	t.Parallel()

	pool, err := NewPool(context.Background(), Config{
		DatabaseURL: "postgres://user:pass@127.0.0.1:1/dbname?sslmode=disable",
	})

	require.Error(t, err)
	require.Nil(t, pool)
	require.ErrorContains(t, err, "ping database")
}
