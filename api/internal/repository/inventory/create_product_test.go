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

func Test_impl_CreateProduct(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath string
		givenCtx     context.Context
		givenProduct model.Product
		mockIDErr    error
		expErr       error
	}

	tcs := map[string]arg{
		"success": {
			givenCtx: context.Background(),
			givenProduct: model.Product{
				Name:        "TestProduct",
				Description: "TestProduct",
				Price:       2500,
				Stock:       5,
				Status:      model.ProductStatusActive,
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			givenProduct: model.Product{
				Name:        "TestProduct",
				Description: "TestProduct",
				Price:       2500,
				Stock:       5,
				Status:      model.ProductStatusActive,
			},
			expErr: context.Canceled,
		},
		"product_constraint": {
			testDataPath: "testdata/success.sql",
			givenCtx:     context.Background(),
			givenProduct: model.Product{
				Name:        "",
				Description: "",
				Price:       2500,
				Stock:       5,
				Status:      model.ProductStatusActive,
			},
			expErr: errors.New("orm: unable to insert into products: pq: new row for relation \"products\" violates check constraint \"products_description_check\""),
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
				createProduct, err := repo.CreateProduct(tc.givenCtx, tc.givenProduct)

				// Then:
				if tc.expErr != nil {
					require.Error(t, err)
					if desc == "product_constraint" {
						// For database constraint errors, just check that the error contains the expected substring
						require.Contains(t, err.Error(), tc.expErr.Error())
					} else {
						require.Equal(t, tc.expErr, pkgerrors.Cause(err))
					}
				} else {
					require.NoError(t, err)
					require.NotEmpty(t, createProduct.ID)
					testutil.Compare(t, tc.givenProduct, createProduct, model.Product{}, "ID", "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
