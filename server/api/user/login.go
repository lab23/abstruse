package user

import (
	"net/http"

	"github.com/lab23/abstruse/internal/auth"
	"github.com/lab23/abstruse/pkg/lib"
	"github.com/lab23/abstruse/server/api/render"
	"github.com/lab23/abstruse/server/core"
)

// HandleLogin returns an http.HandlerFunc that writes JSON encoded
// login data to the http response body.
func HandleLogin(users core.UserStore) http.HandlerFunc {
	type form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type resp struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var f form
		defer r.Body.Close()

		if err := lib.DecodeJSON(r.Body, &f); err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		if users.Login(f.Email, f.Password) {
			user, _ := users.FindEmailOrLogin(f.Email)
			token, err := auth.JWT.CreateJWT(user.Claims())
			if err != nil {
				render.InternalServerError(w, err.Error())
				return
			}
			render.JSON(w, http.StatusOK, resp{Token: token})
			return
		}

		render.UnathorizedError(w, "invalid credentials")
	}
}
