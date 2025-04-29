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

func Test_impl_GetByID(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath string
		givenCtx     context.Context
		givenID      int64
		expUser      model.User
		mockIDErr    error
		expErr       error
	}

	tcs := map[string]arg{
		"success": {
			testDataPath: "testdata/success_get_user.sql",
			givenCtx:     context.Background(),
			givenID:      14753001,
			expUser: model.User{
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive,
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			givenID:  14753001,
			expErr:   context.Canceled,
		},
		"user_not_found": {
			givenCtx: context.Background(),
			givenID:  147530012,
			expErr:   ErrNotFound,
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
				user, err := repo.GetByID(tc.givenCtx, tc.givenID)

				// Then:
				if tc.expErr != nil {
					require.Error(t, err)
					if desc == "duplicate_email" {
						require.Contains(t, err.Error(), tc.expErr.Error())
					} else {
						require.Equal(t, tc.expErr, pkgerrors.Cause(err))
					}
					//require.Equal(t, model.User{}, createdUser)
				} else {
					require.NoError(t, err)
					require.NotEmpty(t, user.ID)
					testutil.Compare(t, tc.expUser, user, model.User{}, "ID", "Password", "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
