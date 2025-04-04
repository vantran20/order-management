package httpserv

import (
	"encoding/json"
	"io"
	"net/http"
)

// ParseJSON parses JSON from the given input
func ParseJSON(r io.ReadCloser, result interface{}) *Error {
	reqBytes, err := io.ReadAll(r)
	defer r.Close()
	if err != nil {
		return &Error{Status: http.StatusBadRequest, Code: "read_body_failed", Desc: err.Error()}
	}

	if err = json.Unmarshal(reqBytes, &result); err != nil {
		return &Error{Status: http.StatusBadRequest, Code: "parse_body_failed", Desc: err.Error()}
	}
	return nil
}
