package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {
	var registerd = []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/static/*", "GET"},
	}

	var app application
	mux := app.routes()
	chiRoutes := mux.(chi.Routes)

	for _, route := range registerd {
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})
	return found
}
