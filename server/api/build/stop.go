package build

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/lab23/abstruse/pkg/lib"
	"github.com/lab23/abstruse/server/api/middlewares"
	"github.com/lab23/abstruse/server/api/render"
	"github.com/lab23/abstruse/server/core"
)

// HandleStop returns an http.HandlerFunc that writes JSON encoded
// result about stopping build to http response body.
func HandleStop(builds core.BuildStore, repos core.RepositoryStore, scheduler core.Scheduler) http.HandlerFunc {
	type form struct {
		ID uint `json:"id" valid:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.ClaimsFromCtx(r.Context())
		var f form
		var err error
		defer r.Body.Close()

		if err = lib.DecodeJSON(r.Body, &f); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		if valid, err := govalidator.ValidateStruct(f); err != nil || !valid {
			render.BadRequestError(w, err.Error())
			return
		}

		build, err := builds.Find(uint(f.ID))
		if err != nil {
			render.NotFoundError(w, err.Error())
			return
		}

		if perms := repos.GetPermissions(build.RepositoryID, claims.ID); !perms.Exec {
			render.UnathorizedError(w, "permission denied")
			return
		}

		if err := scheduler.StopBuild(build.ID); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, render.Empty{})
	}
}
