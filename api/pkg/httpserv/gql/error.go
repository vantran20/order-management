package gql

import (
	"context"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"omg/api/pkg/httpserv"
)

const (
	// gqlValidationFailedCode represents the default code for validation failure
	gqlValidationFailedCode = "gql_validation_failed"
	// gqlParseFailedCode indicates that the parsing of gql failed
	gqlParseFailedCode = "gql_parse_failed"
)

// errorPresenter converts Error or error to a presentable format and also handles reporting of errors.
func errorPresenter(isIntrospectionEnabled bool) func(ctx context.Context, err error) *gqlerror.Error {
	return func(ctx context.Context, err error) *gqlerror.Error {
		if err == nil {
			return nil
		}

		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		var werr *httpserv.Error

		if !errors.As(err, &werr) {
			if gqlErr.Extensions != nil { // Possibly a gqlgen parsed error
				werr = parseGQLError(gqlErr, isIntrospectionEnabled)
			} else { // If unknown error
				werr = httpserv.ErrDefaultInternal
			}
		}

		if werr.Status >= http.StatusInternalServerError && werr.Status != http.StatusServiceUnavailable {
			// Only log internal server (unexpected) errors to reduce noise in logging & sentry
			// For all unexpected errors except service unavailable, generify the message
			werr.Desc = httpserv.DefaultErrorDesc
		}

		if !isIntrospectionEnabled { // We don't want to expose any schema identifiable data when introspection is disabled
			gqlErr.Locations = nil
			gqlErr.Path = nil
		}

		gqlErr.Message = werr.Desc
		gqlErr.Extensions = map[string]interface{}{
			"status":            werr.Status,
			"error":             werr.Code,
			"error_description": werr.Desc,
		}

		return gqlErr
	}
}

// parses the *gqlerror.Error to see if it is gql validation error or parsing failed err and if yes, convert them into
// *httpserv.Error. If introspection is disabled, it will replace the error description with default description
func parseGQLError(gqlErr *gqlerror.Error, isIntrospectionEnabled bool) *httpserv.Error {
	switch gqlErr.Extensions["code"] {
	case errcode.ValidationFailed:
		werr := &httpserv.Error{
			Desc:   gqlErr.Message,
			Code:   gqlValidationFailedCode,
			Status: http.StatusBadRequest,
		}
		if !isIntrospectionEnabled {
			werr.Desc = httpserv.ErrDefaultInternal.Desc
		}
		return werr
	case errcode.ParseFailed:
		werr := &httpserv.Error{
			Desc:   gqlErr.Message,
			Code:   gqlParseFailedCode,
			Status: http.StatusBadRequest,
		}
		if !isIntrospectionEnabled {
			werr.Desc = httpserv.ErrDefaultInternal.Desc
		}
		return werr
	default:
		// If unknown error
		return httpserv.ErrDefaultInternal
	}
}
