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

func Test_impl_Update(t *testing.T) {
	cancelledCtx, c := context.WithCancel(context.Background())
	c()

	type arg struct {
		testDataPath string
		givenCtx     context.Context
		givenUser    model.User
		mockIDErr    error
		expErr       error
	}

	tcs := map[string]arg{
		"success": {
			testDataPath: "testdata/success_get_user.sql",
			givenCtx:     context.Background(),
			givenUser: model.User{
				ID:       14753001,
				Name:     "Test User1",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive,
			},
		},
		"ctx_cancelled": {
			givenCtx: cancelledCtx,
			givenUser: model.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive,
			},
			expErr: context.Canceled,
		},
		"not_found": {
			testDataPath: "testdata/success_get_user.sql",
			givenCtx:     context.Background(),
			givenUser: model.User{
				ID:       14753012,
				Name:     "Test User1",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive,
			},
			expErr: ErrNotFound,
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
				err := repo.Update(tc.givenCtx, tc.givenUser)

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

					updatedUser, er := repo.GetByID(tc.givenCtx, tc.givenUser.ID)
					require.NoError(t, er)
					testutil.Compare(t, tc.givenUser, updatedUser, model.User{}, "ID", "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
