package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name             string
		url              string
		expectStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/cat", http.StatusNotFound},
	}

	var app application
	routes := app.routes()

	// 启动一个测试服务器
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// NOTE: 测试时修正模板地址
	pathToTemplate = "./../../templates/"

	for _, e := range theTests {
		// ts.URL -> base URL of form http://ipaddr:port with no trailing slash
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectStatusCode, resp.StatusCode)
		}
	}
}
