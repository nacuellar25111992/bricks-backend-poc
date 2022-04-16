package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/api/env", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()

	handler := http.HandlerFunc(srv.envHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	testResult := []string{}
	_ = json.Unmarshal(rr.Body.Bytes(), &testResult)
	require.NotEmpty(t, testResult)

	pathEnvFound := false
	for _, envKeyValue := range testResult {
		envKey := strings.Split(envKeyValue, "=")[0]
		if envKey == "PATH" {
			pathEnvFound = true
			break
		}
	}
	if !pathEnvFound {
		require.Fail(t, "path env var not found")
	}
}
