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

func TestTokenHandler(t *testing.T) {

	req, err := http.NewRequest("POST", "/token", strings.NewReader("test-user"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()

	handler := http.HandlerFunc(srv.tokenGenerateHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var token TokenResponse

	if err := json.Unmarshal(rr.Body.Bytes(), &token); err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, token.Token)
}
