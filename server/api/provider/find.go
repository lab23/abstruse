package provider

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/lab23/abstruse/server/api/middlewares"
	"github.com/lab23/abstruse/server/api/render"
	"github.com/lab23/abstruse/server/core"
)

// HandleFind writes JSON encoded provider data to the http response body.
func HandleFind(providers core.ProviderStore, users core.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.ClaimsFromCtx(r.Context())
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		user, err := users.Find(claims.ID)
		if err != nil {
			render.UnathorizedError(w, err.Error())
			return
		}

		provider, err := providers.Find(uint(id))
		if err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		if provider.UserID == claims.ID || user.Role == "admin" {
			render.JSON(w, http.StatusOK, provider)
		}

		render.UnathorizedError(w, err.Error())
	}
}
