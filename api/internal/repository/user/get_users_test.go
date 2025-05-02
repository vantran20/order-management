package user

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

func Test_impl_GetUsers(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath string
		givenCtx     context.Context
		expUser      []model.User
		mockIDErr    error
		expErr       error
	}

	tcs := map[string]arg{
		"success": {
			testDataPath: "testdata/success_get_user.sql",
			givenCtx:     context.Background(),
			expUser: []model.User{
				{
					Name:   "Test User",
					Email:  "test@example.com",
					Status: model.UserStatusActive,
				},
				{
					Name:   "Test User2",
					Email:  "test2@example.com",
					Status: model.UserStatusActive,
				},
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			expErr:   context.Canceled,
		},
		"empty_list": {
			givenCtx: context.Background(),
			expUser:  nil,
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
				users, err := repo.GetUsers(tc.givenCtx)

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
					testutil.Compare(t, tc.expUser, users, model.User{}, "ID", "Password", "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
