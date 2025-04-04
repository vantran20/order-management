package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

// Handler creates a http.Handler from the given ExecutableSchema and returns it
func Handler(es graphql.ExecutableSchema, isIntrospectionEnabled bool) http.Handler {
	srv := handler.New(es)
	srv.AddTransport(transport.POST{})
	srv.SetErrorPresenter(errorPresenter(isIntrospectionEnabled))
	if isIntrospectionEnabled {
		srv.Use(extension.Introspection{})
	}
	return srv
}
