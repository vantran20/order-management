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

func Test_impl_CreateOrderItem(t *testing.T) {
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
			testDataPath: "testdata/success.sql",
			givenCtx:     context.Background(),
			givenOrderItem: model.OrderItem{
				OrderID:   14753010,
				ProductID: 14753010,
				Quantity:  10,
				Price:     2000,
			},
		},
		"ctx_cancelled": {
			testDataPath: "testdata/success.sql",
			givenCtx:     cancelledCtx,
			givenOrderItem: model.OrderItem{
				OrderID:   14753010,
				ProductID: 14753010,
				Quantity:  10,
				Price:     2000,
			},
			expErr: context.Canceled,
		},
		"order_not_found": {
			testDataPath: "testdata/success.sql",
			givenCtx:     context.Background(),
			givenOrderItem: model.OrderItem{
				OrderID:   14753020,
				ProductID: 14753010,
				Quantity:  10,
				Price:     2000,
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
				createdOrderItem, err := repo.CreateOrderItem(tc.givenCtx, tc.givenOrderItem)

				// Then:
				if tc.expErr != nil {
					require.Error(t, err)
					if desc == "order_not_found" {
						// For database constraint errors, just check that the error contains the expected substring
						require.Contains(t, err.Error(), tc.expErr.Error())
					} else {
						require.Equal(t, tc.expErr, pkgerrors.Cause(err))
					}
				} else {
					require.NoError(t, err)
					require.NotEmpty(t, createdOrderItem.ID)
					testutil.Compare(t, tc.givenOrderItem, createdOrderItem, model.OrderItem{}, "ID", "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
