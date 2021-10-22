package config

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersionHandler(t *testing.T) {
	expected := "Version: " + version
	req, err := http.NewRequest(http.MethodGet, "/version", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(getVersion)
	handler.ServeHTTP(res, req)

	if res.Body.String() != expected {
		t.Errorf("Error versionHandler, got: %s, want: %s", res.Body.String(), expected)
	}
}
