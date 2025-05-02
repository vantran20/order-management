package inventory

import (
	"context"
	"errors"
	"testing"

	"omg/api/internal/model"
	"omg/api/internal/repository/generator"
	"omg/api/pkg/db/pg"
	"omg/api/pkg/testutil"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_impl_CreateOrder(t *testing.T) {
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
			testDataPath: "testdata/success.sql",
			givenCtx:     context.Background(),
			givenOrder: model.Order{
				UserID:    14753001,
				Status:    model.OrderStatusPending,
				TotalCost: 10,
			},
		},
		"ctx_cancelled": {
			testDataPath: "testdata/success.sql",
			givenCtx:     cancelledCtx,
			givenOrder: model.Order{
				UserID:    14753001,
				Status:    model.OrderStatusPending,
				TotalCost: 10,
			},
			expErr: context.Canceled,
		},
		"user_not_found": {
			testDataPath: "testdata/success.sql",
			givenCtx:     context.Background(),
			givenOrder: model.Order{
				ID:        14753010,
				UserID:    14753011,
				Status:    model.OrderStatusPending,
				TotalCost: 10,
			},
			expErr: errors.New("foreign key constraint"),
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
				createdOrder, err := repo.CreateOrder(tc.givenCtx, tc.givenOrder)

				// Then:
				if tc.expErr != nil {
					require.Error(t, err)
					if desc == "user_not_found" {
						// For database constraint errors, just check that the error contains the expected substring
						require.Contains(t, err.Error(), tc.expErr.Error())
					} else {
						require.Equal(t, tc.expErr, pkgerrors.Cause(err))
					}
				} else {
					require.NoError(t, err)
					require.NotEmpty(t, createdOrder.ID)
					testutil.Compare(t, tc.givenOrder, createdOrder, model.Order{}, "ID", "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
