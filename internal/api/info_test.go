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

func TestInfoHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/api/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()

	handler := http.HandlerFunc(srv.infoHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	// TODO: replace this with an object assertion please.

	expectedBodyString := strings.ReplaceAll(
		`
	{
		"hostname": "localhost",
		"version": "0.1.0",
		"revision": "unknown",
		"message": "",
		"goos": "darwin",
		"goarch": "arm64",
		"runtime": "go1.18",
		"num_goroutine": "2",
		"num_cpu": "8"
	}
	`, "\t", "")
	expectedBodyString = strings.ReplaceAll(expectedBodyString, "\n", "")
	expectedBodyString = strings.ReplaceAll(expectedBodyString, " ", "")
	expectedBodyPrettified, _ := json.MarshalIndent(expectedBodyString, "", "    ")
	testResultBodyString := strings.ReplaceAll(rr.Body.String(), "\t", "")
	testResultBodyString = strings.ReplaceAll(testResultBodyString, "\n", "")
	testResultBodyString = strings.ReplaceAll(testResultBodyString, " ", "")
	testResultBodyPrettified, _ := json.MarshalIndent(testResultBodyString, "", "    ")
	assert.Equal(t, string(expectedBodyPrettified), string(testResultBodyPrettified))
}
