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

func TestImpl_Create(t *testing.T) {
	type arg struct {
		givenInput        model.CreateUserInput
		mockUserRepoOut   model.User
		mockUserRepoErr   error
		mockGetByEmailErr error
		mockCreateUserErr error
		expRepoMockCalled bool
		expResult         model.User
		expErr            error
	}

	tcs := map[string]arg{
		"success": {
			givenInput: model.CreateUserInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockUserRepoOut: model.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Status:   model.UserStatusActive,
				Password: "$2a$10$somehashedpassword",
			},
			mockGetByEmailErr: user.ErrNotFound,
			expRepoMockCalled: true,
			expResult: model.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Status:   model.UserStatusActive,
				Password: "$2a$10$somehashedpassword",
			},
			expErr: nil,
		},
		"user_already_exists": {
			givenInput: model.CreateUserInput{
				Name:     "Test User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockUserRepoOut: model.User{
				ID:     1,
				Email:  "existing@example.com",
				Status: model.UserStatusActive,
			},
			mockGetByEmailErr: nil,
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            ErrUserAlreadyExists,
		},
		"get_by_email_error": {
			givenInput: model.CreateUserInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockGetByEmailErr: errors.New("database error"),
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            errors.New("database error"),
		},
		"create_user_error": {
			givenInput: model.CreateUserInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockGetByEmailErr: user.ErrNotFound,
			mockCreateUserErr: errors.New("failed to create user"),
			expRepoMockCalled: true,
			expResult:         model.User{},
			expErr:            errors.New("failed to create user"),
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			userRepo := user.MockRepository{}

			if tc.expRepoMockCalled {
				// Mock GetByEmail call
				userRepo.On("GetByEmail", mock.Anything, tc.givenInput.Email).Return(model.User{}, tc.mockGetByEmailErr)

				// If we expect to reach the CreateUser call (when user doesn't exist)
				if errors.Is(tc.mockGetByEmailErr, user.ErrNotFound) {
					// We need to match any User model here since the password will be dynamically hashed
					userRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u model.User) bool {
						return u.Email == tc.givenInput.Email &&
							u.Name == tc.givenInput.Name &&
							u.Status == model.UserStatusActive &&
							len(u.Password) > 0 // Just check that password is not empty
					})).Return(tc.mockUserRepoOut, tc.mockCreateUserErr)
				}
			}

			repo := repository.MockRegistry{}
			repo.On("User").Return(&userRepo)

			impl := impl{repo: &repo}

			// When:
			result, err := impl.Create(context.Background(), tc.givenInput)

			// Then:
			if tc.expErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expErr.Error(), err.Error())
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
