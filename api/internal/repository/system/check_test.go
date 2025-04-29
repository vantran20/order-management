package system

import (
	"context"
	"testing"

	"omg/api/pkg/db/pg"
	"omg/api/pkg/testutil"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestImpl_CheckDB(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()
	type arg struct {
		givenCtx context.Context
		expErr   error
	}
	tcs := map[string]arg{
		"success": {givenCtx: context.Background()},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			expErr:   context.Canceled,
		},
	}
	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			testutil.WithTxDB(t, func(dbConn pg.BeginnerExecutor) {
				// Given:
				repo := New(dbConn)

				// When:
				err := repo.CheckDB(tc.givenCtx)

				// Then:
				require.Equal(t, tc.expErr, pkgerrors.Cause(err))
			})
		})
	}
}
