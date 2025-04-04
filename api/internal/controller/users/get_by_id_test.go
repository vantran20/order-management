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

func TestImpl_GetByID(t *testing.T) {
	type arg struct {
		givenID           int64
		mockGetByIDOut    model.User
		mockGetByIDErr    error
		expRepoMockCalled bool
		expResult         model.User
		expErr            error
	}

	tcs := map[string]arg{
		"success": {
			givenID: 1,
			mockGetByIDOut: model.User{
				ID:     1,
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive,
			},
			mockGetByIDErr:    nil,
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
			givenID:           999,
			mockGetByIDErr:    user.ErrNotFound,
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            ErrUserNotFound,
		},
		"database_error": {
			givenID:           1,
			mockGetByIDErr:    errors.New("database error"),
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
				// Mock GetByID call
				userRepo.On("GetByID", mock.Anything, tc.givenID).Return(tc.mockGetByIDOut, tc.mockGetByIDErr)
			}

			repo := repository.MockRegistry{}
			repo.On("User").Return(&userRepo)

			impl := impl{repo: &repo}

			// When:
			result, err := impl.GetByID(context.Background(), tc.givenID)

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
