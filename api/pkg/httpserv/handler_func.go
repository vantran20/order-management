package httpserv

import (
	"net/http"
)

// ErrHandlerFunc is a convenience wrapper for http.HandlerFunc that handles error transformation and reporting.
// This handler results in a more natural coding style for golang which is to return the error from the func.
// If the error is not a httpserv.Error with status < 5xx & 503, then it will report the error.
func ErrHandlerFunc(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := h(w, r); err != nil {
			RespondJSON(ctx, w, err)

			if werr, ok := err.(*Error); ok {
				if werr.Status < http.StatusInternalServerError || werr.Status == http.StatusServiceUnavailable {
					return
				}
			}
		}
	}
}
