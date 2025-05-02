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

func Test_impl_ListProducts(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath string
		givenCtx     context.Context
		expProducts  []model.Product
		mockIDErr    error
		expErr       error
	}

	tcs := map[string]arg{
		"success": {
			testDataPath: "testdata/success_get_data.sql",
			givenCtx:     context.Background(),
			expProducts: []model.Product{
				{
					ID:          14753010,
					Name:        "Test Product",
					Description: "test",
					Status:      model.ProductStatusActive,
					Price:       2000,
					Stock:       100,
				},
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			expErr:   context.Canceled,
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
				products, err := repo.ListProducts(tc.givenCtx)

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
					testutil.Compare(t, tc.expProducts, products, model.Product{}, "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
