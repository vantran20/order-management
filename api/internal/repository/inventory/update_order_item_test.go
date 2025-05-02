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

func Test_impl_UpdateOrderItem(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath   string
		givenCtx       context.Context
		givenOrderItem model.OrderItem
		mockIDErr      error
		expErr         error
	}

	tcs := map[string]arg{
		"success": {
			testDataPath: "testdata/success_get_data.sql",
			givenCtx:     context.Background(),
			givenOrderItem: model.OrderItem{
				ID:        14753001,
				OrderID:   14753010,
				ProductID: 14753010,
				Quantity:  30,
				Price:     60000,
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			givenOrderItem: model.OrderItem{
				ID:        14753001,
				OrderID:   14753010,
				ProductID: 14753010,
				Quantity:  30,
				Price:     60000,
			},
			expErr: context.Canceled,
		},
		"not_found": {
			testDataPath: "testdata/success_get_data.sql",
			givenCtx:     context.Background(),
			givenOrderItem: model.OrderItem{
				ID:        14753012,
				OrderID:   14753010,
				ProductID: 14753010,
				Quantity:  30,
				Price:     60000,
			},
			expErr: ErrOrderItemNotFound,
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
				orderItemUpdated, err := repo.UpdateOrderItem(tc.givenCtx, tc.givenOrderItem)

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

					testutil.Compare(t, tc.givenOrderItem, orderItemUpdated, model.OrderItem{}, "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
