package repo

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/lab23/abstruse/pkg/lib"
	"github.com/lab23/abstruse/server/api/middlewares"
	"github.com/lab23/abstruse/server/api/render"
	"github.com/lab23/abstruse/server/core"
)

// HandleCreateMount returns an http.HandlerFunc that writes json encoded
// result about creating env variable to the http response body.
func HandleCreateMount(mount core.MountsStore, repos core.RepositoryStore) http.HandlerFunc {
	type form struct {
		Host      string `json:"host" valid:"required"`
		Container string `json:"container" valid:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.ClaimsFromCtx(r.Context())
		var f form
		var err error
		defer r.Body.Close()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		if err = lib.DecodeJSON(r.Body, &f); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		if valid, err := govalidator.ValidateStruct(f); err != nil || !valid {
			render.BadRequestError(w, err.Error())
			return
		}

		if perm := repos.GetPermissions(uint(id), claims.ID); !perm.Write {
			render.UnathorizedError(w, "permission denied")
			return
		}

		mnt := &core.Mount{
			Host:         f.Host,
			Container:    f.Container,
			RepositoryID: uint(id),
		}

		if err := mount.Create(mnt); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}
		render.JSON(w, http.StatusOK, mnt)
	}
}
