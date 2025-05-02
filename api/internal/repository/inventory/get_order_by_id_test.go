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

func Test_impl_GetOrderByID(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath string
		givenCtx     context.Context
		givenID      int64
		expOrder     model.Order
		mockIDErr    error
		expErr       error
	}

	tcs := map[string]arg{
		"success": {
			testDataPath: "testdata/success_get_data.sql",
			givenCtx:     context.Background(),
			givenID:      14753010,
			expOrder: model.Order{
				ID:        14753010,
				UserID:    14753001,
				Status:    model.OrderStatusPending,
				TotalCost: 20,
				OrderItems: []model.OrderItem{
					{
						ID:        14753001,
						OrderID:   14753010,
						ProductID: 14753010,
						Quantity:  20,
						Price:     2000,
					},
				},
			},
		},
		"success_with_empty_items": {
			testDataPath: "testdata/success_get_data.sql",
			givenCtx:     context.Background(),
			givenID:      14753011,
			expOrder: model.Order{
				ID:        14753011,
				UserID:    14753001,
				Status:    model.OrderStatusPending,
				TotalCost: 10,
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			givenID:  14753001,
			expErr:   context.Canceled,
		},
		"order_not_found": {
			givenCtx: context.Background(),
			givenID:  147530012,
			expErr:   ErrOrderNotFound,
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
				order, err := repo.GetOrderByID(tc.givenCtx, tc.givenID)

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
					require.NotEmpty(t, order.ID)
					testutil.Compare(t, tc.expOrder, order, model.Order{}, "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
