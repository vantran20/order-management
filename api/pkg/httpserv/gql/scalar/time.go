package scalar

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"omg/api/pkg/httpserv"
)

// MarshalTime is custom impl that overwrites gqlgen's time.RFC3339Nano format to time.RFC3339
func MarshalTime(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(t.Format(time.RFC3339)))
	})
}

// UnmarshalTime is custom impl that overwrites gqlgen's time.RFC3339Nano format to time.RFC3339
func UnmarshalTime(v interface{}) (time.Time, error) {
	var err error
	switch v := v.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, &httpserv.Error{
				Status: http.StatusBadRequest,
				Code:   "invalid_gql_field_value",
				Desc:   "Unable to convert value to time",
			}
		}
		return t, nil
	default:
		err = &httpserv.Error{
			Status: http.StatusBadRequest,
			Code:   "invalid_gql_field_type",
			Desc:   "time should be RFC3339 formatted string",
		}
	}
	return time.Time{}, err
}
