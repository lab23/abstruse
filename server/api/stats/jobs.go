package stats

import (
	"net/http"
	"time"

	"github.com/lab23/abstruse/server/api/render"
	"github.com/lab23/abstruse/server/core"
)

// HandleJobs returns an http.HandlerFunc that writes JSON encoded
// result about jobs statistics to the http response body.
func HandleJobs(jobs core.JobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		layout := "2006-01-02"
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		tfrom, err := time.Parse(layout, from)
		if err != nil {
			tfrom = time.Now().AddDate(0, 0, -7)
		}
		tto, err := time.Parse(layout, to)
		if err != nil {
			tto = time.Now()
		}

		tfrom = time.Date(tfrom.Year(), tfrom.Month(), tfrom.Day(), 0, 0, 0, 0, time.UTC)
		tto = time.Date(tto.Year(), tto.Month(), tto.Day(), 23, 59, 59, 0, time.UTC)

		jobs, err := jobs.List(tfrom, tto)
		if err != nil {
			render.InternalServerError(w, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, jobs)
	}
}
