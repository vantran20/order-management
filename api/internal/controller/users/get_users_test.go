package users

import (
	"context"
	"errors"
	"testing"

	"omg/api/internal/model"
	"omg/api/internal/repository"
	"omg/api/internal/repository/user"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_impl_GetUsers(t *testing.T) {
	type arg struct {
		mockGetUsersOut   []model.User
		mockGetUsersErr   error
		expRepoMockCalled bool
		expResult         []model.User
		expErr            error
	}

	tcs := map[string]arg{
		"success": {
			mockGetUsersOut: []model.User{
				{
					ID:     1,
					Name:   "Test User",
					Email:  "test@example.com",
					Status: model.UserStatusActive,
				},
			},
			mockGetUsersErr:   nil,
			expRepoMockCalled: true,
			expResult: []model.User{
				{
					ID:     1,
					Name:   "Test User",
					Email:  "test@example.com",
					Status: model.UserStatusActive,
				},
			},
			expErr: nil,
		},
		"empty": {
			mockGetUsersOut:   nil,
			mockGetUsersErr:   nil,
			expRepoMockCalled: true,
			expResult:         nil,
			expErr:            nil,
		},
		"database_error": {
			mockGetUsersErr:   errors.New("database error"),
			expRepoMockCalled: true,
			expResult:         nil,
			expErr:            errors.New("database error"),
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			userRepo := user.MockRepository{}

			if tc.expRepoMockCalled {
				// Mock GetUsers call
				userRepo.On("GetUsers", mock.Anything).Return(tc.mockGetUsersOut, tc.mockGetUsersErr)
			}

			repo := repository.MockRegistry{}
			repo.On("User").Return(&userRepo)

			impl := impl{repo: &repo}

			// When:
			result, err := impl.GetUsers(context.Background())

			// Then:
			if tc.expErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expErr.Error(), err.Error())
				require.Equal(t, tc.expResult, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}

			userRepo.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
