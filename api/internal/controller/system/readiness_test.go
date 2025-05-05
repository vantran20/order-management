package system

import (
	"context"
	"errors"
	"testing"

	"omg/api/internal/repository"
	"omg/api/internal/repository/system"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestImpl_CheckReadiness(t *testing.T) {
	type arg struct {
		givenIAMReady    bool
		mockDBRepoOutErr error
		expDBMockCalled  bool
		expErr           error
	}
	tcs := map[string]arg{
		"success": {
			givenIAMReady:   true,
			expDBMockCalled: true,
		},
		"dberr": {
			givenIAMReady:    true,
			expDBMockCalled:  true,
			mockDBRepoOutErr: errors.New("some error"),
			expErr:           errors.New("some error"),
		},
	}
	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			systemRepo := system.MockRepository{}
			if tc.expDBMockCalled {
				systemRepo.ExpectedCalls = []*mock.Call{
					systemRepo.On("CheckDB", mock.Anything).Return(tc.mockDBRepoOutErr),
				}
			}

			repo := repository.MockRegistry{}
			if tc.expDBMockCalled {
				repo.ExpectedCalls = []*mock.Call{
					repo.On("System").Return(&systemRepo),
				}
			}

			c := New(&repo)

			// When:
			err := c.CheckReadiness(context.Background())

			// Then:
			require.Equal(t, tc.expErr, pkgerrors.Cause(err))
		})
	}
}
