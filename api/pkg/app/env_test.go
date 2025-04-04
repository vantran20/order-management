package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnv_String(t *testing.T) {
	require.Equal(t, "prod", EnvProd.String())
	require.Equal(t, "qa", EnvQA.String())
	require.Equal(t, "dev", EnvDev.String())
	require.Equal(t, "local", EnvLocal.String())
	require.Equal(t, "test", EnvTest.String())
	require.Equal(t, "abcd", Env("abcd").String())
	require.Equal(t, "", Env("").String())
}

func TestEnv_Valid(t *testing.T) {
	require.True(t, EnvProd.Valid())
	require.True(t, EnvQA.Valid())
	require.True(t, EnvDev.Valid())
	require.True(t, EnvLocal.Valid())
	require.True(t, EnvTest.Valid())
	require.False(t, Env("abcd").Valid())
	require.False(t, Env("").Valid())
}
