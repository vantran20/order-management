package httpserv

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrHandlerFunc(t *testing.T) {
	type arg struct {
		givenErr  error
		expStatus int
		expErr    *Error
	}
	tcs := map[string]arg{
		"no_error": {
			expStatus: http.StatusOK,
		},
		"go_error": {
			givenErr:  errors.New("some error"),
			expStatus: http.StatusInternalServerError,
			expErr:    ErrDefaultInternal,
		},
		"web_error_500": {
			givenErr:  &Error{Status: http.StatusInternalServerError, Code: "code", Desc: "desc"},
			expStatus: http.StatusInternalServerError,
			expErr:    &Error{Status: http.StatusInternalServerError, Code: "code", Desc: DefaultErrorDesc},
		},
		"web_error_503": {
			givenErr:  &Error{Status: http.StatusServiceUnavailable, Code: "code", Desc: "desc"},
			expStatus: http.StatusServiceUnavailable,
			expErr:    &Error{Code: "code", Desc: "desc"},
		},
		"web_error_400": {
			givenErr:  &Error{Status: http.StatusBadRequest, Code: "code", Desc: "desc"},
			expStatus: http.StatusBadRequest,
			expErr:    &Error{Code: "code", Desc: "desc"},
		},
	}
	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// When:
			ErrHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				return tc.givenErr
			}).ServeHTTP(w, r)

			// Then:
			require.Equal(t, tc.expStatus, w.Code)
			var actErr Error
			err := ParseJSON(w.Result().Body, &actErr)
			if tc.expErr == nil {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expErr.Code, actErr.Code)
				require.Equal(t, tc.expErr.Desc, actErr.Desc)
			}
		})
	}
}
