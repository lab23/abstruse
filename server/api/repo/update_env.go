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

// HandleUpdateEnv returns an http.HandlerFunc that writes json encoded
// result about updating env variable to the http response body.
func HandleUpdateEnv(envVariables core.EnvVariableStore, repos core.RepositoryStore) http.HandlerFunc {
	type form struct {
		ID     uint   `json:"id" valid:"required"`
		Key    string `json:"key" valid:"required"`
		Value  string `json:"value"`
		Secret bool   `json:"secret"`
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

		env, err := envVariables.Find(f.ID)
		if err != nil {
			render.NotFoundError(w, err.Error())
			return
		}

		env.Key = f.Key
		env.Secret = f.Secret
		if f.Value != "" {
			env.Value = f.Value
		}

		if err := envVariables.Update(env); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, env)
	}
}
