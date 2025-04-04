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

func TestImpl_Update(t *testing.T) {
	type arg struct {
		givenInput        model.UpdateUserInput
		mockGetByIDOut    model.User
		mockGetByIDErr    error
		mockUpdateErr     error
		expRepoMockCalled bool
		expResult         model.User
		expErr            error
	}

	tcs := map[string]arg{
		"success": {
			givenInput: model.UpdateUserInput{
				ID:       1,
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive,
			},
			mockGetByIDOut: model.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Status:   model.UserStatusActive,
				Password: "$2a$10$oldhashpassword",
			},
			mockGetByIDErr:    nil,
			mockUpdateErr:     nil,
			expRepoMockCalled: true,
			expResult: model.User{
				ID:       1,
				Name:     "Updated User",
				Email:    "updated@example.com",
				Status:   model.UserStatusActive,
				Password: "$2a$10$newhashpassword",
			},
			expErr: nil,
		},
		"user_not_found": {
			givenInput: model.UpdateUserInput{
				ID:       999,
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
			},
			mockGetByIDErr:    user.ErrNotFound,
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            ErrUserNotFound,
		},
		"get_user_error": {
			givenInput: model.UpdateUserInput{
				ID:       1,
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
			},
			mockGetByIDErr:    errors.New("database error"),
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            errors.New("database error"),
		},
		"update_error": {
			givenInput: model.UpdateUserInput{
				ID:       1,
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive,
			},
			mockGetByIDOut: model.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Status:   model.UserStatusActive,
				Password: "$2a$10$oldhashpassword",
			},
			mockGetByIDErr:    nil,
			mockUpdateErr:     errors.New("update error"),
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            errors.New("update error"),
		},
		"update_not_found": {
			givenInput: model.UpdateUserInput{
				ID:       1,
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive,
			},
			mockGetByIDOut: model.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Status:   model.UserStatusActive,
				Password: "$2a$10$oldhashpassword",
			},
			mockGetByIDErr:    nil,
			mockUpdateErr:     user.ErrNotFound,
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            ErrUserNotFound,
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			userRepo := user.MockRepository{}

			if tc.expRepoMockCalled {
				// Mock GetByID call
				userRepo.On("GetByID", mock.Anything, tc.givenInput.ID).Return(tc.mockGetByIDOut, tc.mockGetByIDErr)

				// If we expect to reach the Update call (when user exists)
				if tc.mockGetByIDErr == nil {
					// We need to match any User model here since the password will be dynamically hashed
					userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u model.User) bool {
						return u.ID == tc.givenInput.ID &&
							u.Name == tc.givenInput.Name &&
							u.Email == tc.givenInput.Email &&
							u.Status == tc.givenInput.Status &&
							len(u.Password) > 0 // Just check that password is not empty
					})).Return(tc.mockUpdateErr)
				}
			}

			repo := repository.MockRegistry{}
			repo.On("User").Return(&userRepo)

			impl := impl{repo: &repo}

			// When:
			result, err := impl.Update(context.Background(), tc.givenInput)

			// Then:
			if tc.expErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expErr.Error(), err.Error())
				require.Equal(t, model.User{}, result)
			} else {
				require.NoError(t, err)
				// For the password field, we just check it's not empty since bcrypt generates different hashes each time
				require.Equal(t, tc.expResult.ID, result.ID)
				require.Equal(t, tc.expResult.Name, result.Name)
				require.Equal(t, tc.expResult.Email, result.Email)
				require.Equal(t, tc.expResult.Status, result.Status)
				require.NotEmpty(t, result.Password)
			}

			userRepo.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
