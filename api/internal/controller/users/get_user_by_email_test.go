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

func Test_impl_GetByEmail(t *testing.T) {
	type arg struct {
		givenEmail        string
		mockGetByEmailOut model.User
		mockGetByEmailErr error
		expRepoMockCalled bool
		expResult         model.User
		expErr            error
	}

	tcs := map[string]arg{
		"success": {
			givenEmail: "test@example.com",
			mockGetByEmailOut: model.User{
				ID:     1,
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive,
			},
			mockGetByEmailErr: nil,
			expRepoMockCalled: true,
			expResult: model.User{
				ID:     1,
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive,
			},
			expErr: nil,
		},
		"user_not_found": {
			givenEmail:        "123@example.com",
			mockGetByEmailErr: user.ErrNotFound,
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            ErrUserNotFound,
		},
		"database_error": {
			givenEmail:        "test@example.com",
			mockGetByEmailErr: errors.New("database error"),
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            errors.New("database error"),
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			userRepo := user.MockRepository{}

			if tc.expRepoMockCalled {
				// Mock GetByEmail call
				userRepo.On("GetByEmail", mock.Anything, tc.givenEmail).Return(tc.mockGetByEmailOut, tc.mockGetByEmailErr)
			}

			repo := repository.MockRegistry{}
			repo.On("User").Return(&userRepo)

			impl := impl{repo: &repo}

			// When:
			result, err := impl.GetByEmail(context.Background(), tc.givenEmail)

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
