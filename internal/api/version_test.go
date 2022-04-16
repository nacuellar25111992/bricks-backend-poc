package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()

	handler := http.HandlerFunc(srv.versionHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var token map[string]string

	if err := json.Unmarshal(rr.Body.Bytes(), &token); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "0.1.0", token["version"])
	assert.Equal(t, "unknown", token["commit"])
}
