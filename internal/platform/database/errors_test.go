package database

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsUniqueViolation(t *testing.T) {
	t.Parallel()

	assert.True(t, IsUniqueViolation(&pgconn.PgError{Code: UniqueViolationSQLState}))
	assert.False(t, IsUniqueViolation(&pgconn.PgError{Code: "22000"}))
	assert.False(t, IsUniqueViolation(errors.New("boom")))
}
