package api

import (
	"net/http"

	"github.com/nacuellar25111992/bricks-backend-poc/internal/version"
)

// Version godoc
// @Summary Version
// @Description returns bricks-backend-poc version and git commit hash
// @Tags HTTP API
// @Produce json
// @Router /version [get]
// @Success 200 {object} api.MapResponse
func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {

	result := map[string]string{
		"version": version.VERSION,
		"commit":  version.REVISION,
	}

	s.JSONResponse(w, r, result)
}
