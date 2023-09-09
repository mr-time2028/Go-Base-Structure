package user

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

var routes = []string{
	"/login",
	"/refresh",
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("did not find %s in registred routes", route)
	}
}

func Test_Exists_Routes(t *testing.T) {
	chiRoutes := Routes().(chi.Router)
	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func Test_Routes(t *testing.T) {
	router := chi.NewRouter()

	var tests []struct {
		method string
		path   string
	}

	_ = chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		tests = append(tests, struct {
			method string
			path   string
		}{
			method: method,
			path:   route,
		})
		return nil
	})

	// Run the test cases
	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.path, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("Route %s %s not registered correctly", test.method, test.path)
		}
	}
}
