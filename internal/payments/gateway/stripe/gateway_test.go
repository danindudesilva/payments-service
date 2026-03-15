package stripe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_RequiresSecretKey(t *testing.T) {
	t.Parallel()

	gateway, err := New("")

	require.Error(t, err)
	require.Nil(t, gateway)
	require.Contains(t, err.Error(), "stripe secret key must not be empty")
}

func TestNew_ConstructsGateway(t *testing.T) {
	t.Parallel()

	gateway, err := New("sk_test_123")

	require.NoError(t, err)
	require.NotNil(t, gateway)
	require.NotNil(t, gateway.client)
}
