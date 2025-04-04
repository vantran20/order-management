package httpserv

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// NOTE: The tests in this file have been copied from SPD golang web pkg to verify we are not missing anything.
// Later we should improve the tests in this file to make it meet Athena's testing style.

func TestParseJSONBody(t *testing.T) {
	// Given:
	given := []byte(`{"message":"OK"}`)
	r := httptest.NewRequest(http.MethodPost, "/some/path", bytes.NewReader(given))
	reader := r.Body
	var obj Success

	// When:
	err := ParseJSON(reader, &obj)

	// Then:
	require.Nil(t, err)
	require.Equal(t, "OK", obj.Message)
}

func TestParseJSONBody_ReadError(t *testing.T) {
	// Given:
	reader := errReadCloser{failRead: true}
	var obj Success

	// When:
	err := ParseJSON(reader, &obj)

	// Then:
	require.NotNil(t, err)
	require.Equal(t, http.StatusBadRequest, err.Status)
	require.Equal(t, "read_body_failed", err.Code)
	require.Equal(t, "read error", err.Desc)
}

func TestParseJSONBody_UnmarshalError(t *testing.T) {
	// Given:
	r := httptest.NewRequest(http.MethodPost, "/some/path", bytes.NewReader([]byte("not json")))
	reader := r.Body
	var obj Success

	// When:
	err := ParseJSON(reader, &obj)

	// Then:
	require.NotNil(t, err)
	require.Equal(t, http.StatusBadRequest, err.Status)
	require.Equal(t, "parse_body_failed", err.Code)
}

type errReadCloser struct {
	failRead bool
}

func (r errReadCloser) Read([]byte) (n int, err error) {
	if r.failRead {
		return 0, errors.New("read error")
	}
	return 0, io.EOF
}
func (r errReadCloser) Close() error {
	return errors.New("close error")
}
