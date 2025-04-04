package scalar

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"omg/api/pkg/httpserv"
)

// MarshalFloat64 marshals float64 to string
func MarshalFloat64(t float64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(strconv.FormatFloat(t, 'f', -1, 64)))
	})
}

// UnmarshalFloat64 unmarshals string to float64
func UnmarshalFloat64(v interface{}) (float64, error) {
	var err error
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, &httpserv.Error{
				Status: http.StatusBadRequest,
				Code:   "invalid_gql_field_value",
				Desc:   "Unable to convert value to float64",
			}
		}
		return i, nil
	default:
		err = &httpserv.Error{
			Status: http.StatusBadRequest,
			Code:   "invalid_gql_field_type",
			Desc:   fmt.Sprintf("float64 must be of string type, but got %s instead", reflect.TypeOf(v)),
		}
	}
	return 0, err
}
