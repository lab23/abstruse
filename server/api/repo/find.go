package repo

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lab23/abstruse/server/api/middlewares"
	"github.com/lab23/abstruse/server/api/render"
	"github.com/lab23/abstruse/server/core"
)

// HandleFind returns an http.HandlerFunc that writes JSON encoded
// repository result to the http response body.
func HandleFind(repos core.RepositoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.ClaimsFromCtx(r.Context())

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		repo, err := repos.Find(uint(id), claims.ID)
		if err != nil {
			render.NotFoundError(w, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, repo)
	}
}
