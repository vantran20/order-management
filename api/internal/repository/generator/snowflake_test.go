package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitSnowflakeGenerators(t *testing.T) {
	require.Nil(t, InitSnowflakeGenerators())

	require.NotNil(t, ProductIDSNF)
}
