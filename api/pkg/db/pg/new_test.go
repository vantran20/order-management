package pg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitDBConn(t *testing.T) {
	type input struct {
		dbURL    string
		poolSize int
		maxIdle  int
	}
	type output struct {
		expectErr bool
	}
	tests := []struct {
		name string
		in   input
		out  output
	}{
		{
			name: "OK",
			in: input{
				dbURL:    os.Getenv("DB_URL"),
				poolSize: 10,
				maxIdle:  1,
			},
		},
		{
			name: "db url error",
			in: input{
				dbURL: "random string",
			},
			out: output{expectErr: true},
		},
		{
			name: "conn str db error",
			in: input{
				dbURL:    "postgres://unknowndb:@pg:5432/unknowndb",
				poolSize: 10,
				maxIdle:  1,
			},
			out: output{expectErr: true},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPool(
				tc.in.dbURL,
				tc.in.poolSize,
				tc.in.maxIdle,
			)
			if tc.out.expectErr {
				require.Error(t, err, tc.name)
				return
			}
			require.NoError(t, err, tc.name)
		})
	}
}
