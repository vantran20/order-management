package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"omg/api/internal/authenticate"
	"omg/api/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRouter_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type route struct {
		method string
		path   string
	}

	type arg struct {
		givenRouter    Router
		expectedRoutes []route
	}

	tcs := map[string]arg{
		"default configuration": {
			givenRouter: New(
				context.Background(),
				[]string{"*"},
				false,
				nil,
				nil,
				nil,
				nil,
				authenticate.AuthService{},
				ws.NewHub(),
			),
			expectedRoutes: []route{
				// Public routes
				{method: "POST", path: "/public/users/register"},
				{method: "POST", path: "/public/users/login"},
				{method: "GET", path: "/public/users/ws"},

				// Authenticated routes - Users
				{method: "GET", path: "/authenticated/users/profile"},
				{method: "GET", path: "/authenticated/users/:id"},
				{method: "GET", path: "/authenticated/users/list"},
				{method: "PUT", path: "/authenticated/users/update"},
				{method: "POST", path: "/authenticated/users/delete/:id"},

				// Authenticated routes - Products
				{method: "POST", path: "/authenticated/products/create"},
				{method: "PUT", path: "/authenticated/products/update"},
				{method: "POST", path: "/authenticated/products/delete/:id"},
				{method: "GET", path: "/authenticated/products/:id"},
				{method: "GET", path: "/authenticated/products/list"},

				// Authenticated routes - Orders
				{method: "POST", path: "/authenticated/order/create"},
				{method: "PUT", path: "/authenticated/order/update/:id"},
				{method: "GET", path: "/authenticated/order/ws"},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			var routesFound []route

			// Get the handler
			handler := tc.givenRouter.Handler()

			// Create a test request for each expected route
			for _, expectedRoute := range tc.expectedRoutes {
				req := httptest.NewRequest(expectedRoute.method, expectedRoute.path, nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				// If the route exists, it should not return 404
				if w.Code != http.StatusNotFound {
					routesFound = append(routesFound, expectedRoute)
				}
			}

			// Sort both slices for comparison
			sort.Slice(tc.expectedRoutes, func(i, j int) bool {
				if tc.expectedRoutes[i].method != tc.expectedRoutes[j].method {
					return tc.expectedRoutes[i].method < tc.expectedRoutes[j].method
				}
				return tc.expectedRoutes[i].path < tc.expectedRoutes[j].path
			})

			sort.Slice(routesFound, func(i, j int) bool {
				if routesFound[i].method != routesFound[j].method {
					return routesFound[i].method < routesFound[j].method
				}
				return routesFound[i].path < routesFound[j].path
			})

			// Then:
			require.EqualValues(t, tc.expectedRoutes, routesFound)
		})
	}
}
