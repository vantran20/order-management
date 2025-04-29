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
		givenCtx  context.Context
		givenUser model.User
		mockIDErr error
		expErr    error
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
			givenCtx: context.Background(),
			givenUser: model.User{
				Name:     "Duplicate User",
				Email:    "duplicate@example.com",
				Password: "password123",
				Status:   model.UserStatusActive,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			testutil.WithTxDB(t, func(dbConn pg.BeginnerExecutor) {
				// Given:
				repo := New(dbConn)
				require.Nil(t, generator.InitSnowflakeGenerators())

				// For duplicate email test, first create a user with the same email
				if desc == "duplicate_email" {
					initialUser := model.User{
						Name:     "Initial User",
						Email:    "duplicate@example.com",
						Password: "initialpass",
						Status:   model.UserStatusActive,
					}
					_, err := repo.CreateUser(context.Background(), initialUser)
					require.NoError(t, err)

					// Now the expected error will be a duplicate violation
					tc.expErr = errors.New("pq: duplicate key value violates unique constraint")
				}

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
