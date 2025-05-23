package user

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

func TestImpl_CreateUser(t *testing.T) {
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
			givenCtx: context.Background(),
			givenUser: model.User{
				Name:     "Test User",
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
		"duplicate_email": {
			testDataPath: "testdata/success_get_user.sql",
			givenCtx:     context.Background(),
			givenUser: model.User{
				ID:       14753001,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive,
			},
			expErr: errors.New("pq: duplicate key value violates unique constraint"),
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
				createdUser, err := repo.CreateUser(tc.givenCtx, tc.givenUser)

				// Then:
				if tc.expErr != nil {
					require.Error(t, err)
					if desc == "duplicate_email" {
						require.Contains(t, err.Error(), tc.expErr.Error())
					} else {
						require.Equal(t, tc.expErr, pkgerrors.Cause(err))
					}
					require.Equal(t, model.User{}, createdUser)
				} else {
					require.NoError(t, err)
					require.NotEmpty(t, createdUser.ID)
					testutil.Compare(t, tc.givenUser, createdUser, model.User{}, "ID", "CreatedAt", "UpdatedAt")
				}
			})
		})
	}
}
