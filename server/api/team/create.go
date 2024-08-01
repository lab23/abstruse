package team

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/lab23/abstruse/pkg/lib"
	"github.com/lab23/abstruse/server/api/middlewares"
	"github.com/lab23/abstruse/server/api/render"
	"github.com/lab23/abstruse/server/core"
)

// HandleCreate returns an http.HandlerFunc that writes JSON encoded
// result about creating team to the http response body.
func HandleCreate(teams core.TeamStore, users core.UserStore, permissions core.PermissionStore) http.HandlerFunc {
	type repoPerm struct {
		ID    uint `json:"id"`
		Read  bool `json:"read"`
		Write bool `json:"write"`
		Exec  bool `json:"exec"`
	}

	type form struct {
		Name    string     `json:"name" valid:"required"`
		About   string     `json:"about" valid:"required"`
		Color   string     `json:"color" valid:"required"`
		Members []uint     `json:"members"`
		Repos   []repoPerm `json:"repos"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.ClaimsFromCtx(r.Context())
		var f form
		defer r.Body.Close()

		u, err := users.Find(claims.ID)
		if err != nil || u.Role != "admin" {
			render.UnathorizedError(w, err.Error())
			return
		}

		if err := lib.DecodeJSON(r.Body, &f); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		if valid, err := govalidator.ValidateStruct(f); err != nil || !valid {
			render.BadRequestError(w, err.Error())
			return
		}

		team := &core.Team{
			Name:  f.Name,
			About: f.About,
			Color: f.Color,
		}

		if err := teams.Create(team); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		var members []*core.User
		for _, id := range f.Members {
			if user, err := users.Find(id); err == nil {
				members = append(members, user)
			}
		}

		if err := teams.UpdateUsers(team.ID, members); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		for _, perm := range f.Repos {
			if err := permissions.Create(&core.Permission{
				TeamID:       team.ID,
				RepositoryID: perm.ID,
				Read:         perm.Read,
				Write:        perm.Write,
				Exec:         perm.Exec,
			}); err != nil {
				render.InternalServerError(w, err.Error())
				return
			}
		}

		render.JSON(w, http.StatusOK, team)
	}
}
