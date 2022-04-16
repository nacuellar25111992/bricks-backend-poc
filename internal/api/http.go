package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nacuellar25111992/bricks-backend-poc/internal/version"

	"go.uber.org/zap"
)

func versionMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("X-API-Version", version.VERSION)
		r.Header.Set("X-API-Revision", version.REVISION)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) JSONResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON(body))
}

func (s *Server) JSONResponseCode(w http.ResponseWriter, r *http.Request, result interface{}, responseCode int) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(responseCode)
	w.Write(prettyJSON(body))
}

func (s *Server) ErrorResponse(w http.ResponseWriter, r *http.Request, error string, code int) {

	data := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: error,
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("json marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(prettyJSON(body))
}

func prettyJSON(b []byte) []byte {

	var out bytes.Buffer

	_ = json.Indent(&out, b, "", "    ")

	return out.Bytes()
}
