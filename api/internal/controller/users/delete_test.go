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

func TestImpl_Delete(t *testing.T) {
	type arg struct {
		givenID           int64
		mockGetByIDOut    model.User
		mockGetByIDErr    error
		mockUpdateErr     error
		expRepoMockCalled bool
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
			mockUpdateErr:     nil,
			expRepoMockCalled: true,
			expErr:            nil,
		},
		"user_not_found": {
			givenID:           999,
			mockGetByIDErr:    user.ErrNotFound,
			expRepoMockCalled: true,
			expErr:            ErrUserNotFound,
		},
		"get_user_error": {
			givenID:           1,
			mockGetByIDErr:    errors.New("database error"),
			expRepoMockCalled: true,
			expErr:            errors.New("database error"),
		},
		"update_error": {
			givenID: 1,
			mockGetByIDOut: model.User{
				ID:     1,
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive,
			},
			mockGetByIDErr:    nil,
			mockUpdateErr:     errors.New("update error"),
			expRepoMockCalled: true,
			expErr:            errors.New("update error"),
		},
		"update_not_found": {
			givenID: 1,
			mockGetByIDOut: model.User{
				ID:     1,
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive,
			},
			mockGetByIDErr:    nil,
			mockUpdateErr:     user.ErrNotFound,
			expRepoMockCalled: true,
			expErr:            ErrUserNotFound,
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			userRepo := user.MockRepository{}

			if tc.expRepoMockCalled {
				// Mock GetByID call
				userRepo.On("GetByID", mock.Anything, tc.givenID).Return(tc.mockGetByIDOut, tc.mockGetByIDErr)

				// If we expect to reach the Update call (when user exists)
				if tc.mockGetByIDErr == nil {
					// We need to verify the user status is set to deleted
					expectedUpdatedUser := tc.mockGetByIDOut
					expectedUpdatedUser.Status = model.UserStatusDeleted

					userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u model.User) bool {
						return u.ID == tc.givenID && u.Status == model.UserStatusDeleted
					})).Return(tc.mockUpdateErr)
				}
			}

			repo := repository.MockRegistry{}
			repo.On("User").Return(&userRepo)

			impl := impl{repo: &repo}

			// When:
			err := impl.Delete(context.Background(), tc.givenID)

			// Then:
			if tc.expErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
