package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthzHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()

	handler := http.HandlerFunc(srv.healthzHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusServiceUnavailable, rr.Code)
}

func TestReadyzHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/readyz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()

	handler := http.HandlerFunc(srv.readyzHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusServiceUnavailable, rr.Code)
}
