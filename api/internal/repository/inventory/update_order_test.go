package inventory

import (
	"context"
	"testing"

	"omg/api/internal/model"
	"omg/api/internal/repository/generator"
	"omg/api/pkg/db/pg"
	"omg/api/pkg/testutil"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_impl_UpdateOrder(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath string
		givenCtx     context.Context
		givenOrder   model.Order
		mockIDErr    error
		expErr       error
	}

	tcs := map[string]arg{
		"success": {
			testDataPath: "testdata/success_get_data.sql",
			givenCtx:     context.Background(),
			givenOrder: model.Order{
				ID:        14753010,
				UserID:    14753001,
				Status:    model.OrderStatusPending,
				TotalCost: 30,
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			givenOrder: model.Order{
				ID:        14753010,
				UserID:    14753001,
				Status:    model.OrderStatusPending,
				TotalCost: 30,
			},
			expErr: context.Canceled,
		},
		"not_found": {
			testDataPath: "testdata/success_get_data.sql",
			givenCtx:     context.Background(),
			givenOrder: model.Order{
				ID:        14753012,
				UserID:    14753001,
				Status:    model.OrderStatusPending,
				TotalCost: 30,
			},
			expErr: ErrOrderNotFound,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			testutil.WithTxDB(t, func(dbConn pg.BeginnerExecutor) {
				// Given:
				if tc.testDataPath != "" {
					testutil.LoadTestSQLFile(t, dbConn, tc.testDataPath)
				}

				repo := New(dbConn)
				require.Nil(t, generator.InitSnowflakeGenerators())

				// When:
				orderUpdated, err := repo.UpdateOrder(tc.givenCtx, tc.givenOrder)

				// Then:
				if tc.expErr != nil {
					require.Error(t, err)
					if desc == "duplicate_email" {
						require.Contains(t, err.Error(), tc.expErr.Error())
					} else {
						require.Equal(t, tc.expErr, pkgerrors.Cause(err))
					}
				} else {
					require.NoError(t, err)

					testutil.Compare(t, tc.givenOrder, orderUpdated, model.Order{}, "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
