package pg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExponentialBackOff(t *testing.T) {
	// Given:
	b := ExponentialBackOff(2, 5*time.Second)

	// When & Then:
	require.Equal(t, 500*time.Millisecond, b.NextBackOff())
	require.Equal(t, 750*time.Millisecond, b.NextBackOff())
	require.Equal(t, -1*time.Nanosecond, b.NextBackOff())

	// Given:
	b = ExponentialBackOff(0, 30*time.Second)

	// When & Then:
	require.Equal(t, -1*time.Nanosecond, b.NextBackOff())

	// Given:
	b = ExponentialBackOff(3, -1*time.Second)

	// When & Then:
	require.Equal(t, -1*time.Nanosecond, b.NextBackOff())
}
