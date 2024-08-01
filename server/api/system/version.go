package system

import (
	"net/http"

	"github.com/lab23/abstruse/internal/version"
	"github.com/lab23/abstruse/server/api/render"
)

// HandleVersion returns an http.HandlerFunc that writes JSON
// encoded version data to the http response body.
func HandleVersion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, http.StatusOK, version.GetBuildInfo())
	}
}
