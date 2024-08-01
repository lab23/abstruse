//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/lab23/abstruse/server/api"
	"github.com/lab23/abstruse/server/http"
	"github.com/lab23/abstruse/server/logger"
	"github.com/lab23/abstruse/server/scheduler"
	"github.com/lab23/abstruse/server/service/stats"
	"github.com/lab23/abstruse/server/store"
	"github.com/lab23/abstruse/server/store/build"
	"github.com/lab23/abstruse/server/store/envvariable"
	"github.com/lab23/abstruse/server/store/job"
	"github.com/lab23/abstruse/server/store/mounts"
	"github.com/lab23/abstruse/server/store/permission"
	"github.com/lab23/abstruse/server/store/provider"
	"github.com/lab23/abstruse/server/store/repo"
	"github.com/lab23/abstruse/server/store/team"
	"github.com/lab23/abstruse/server/store/user"
	"github.com/lab23/abstruse/server/worker"
	"github.com/lab23/abstruse/server/ws"
)

func CreateApp() (*app, error) {
	panic(wire.Build(wire.NewSet(
		wire.NewSet(store.New),
		wire.NewSet(api.New),
		wire.NewSet(user.New),
		wire.NewSet(team.New),
		wire.NewSet(permission.New),
		wire.NewSet(provider.New),
		wire.NewSet(build.New),
		wire.NewSet(job.New),
		wire.NewSet(repo.New),
		wire.NewSet(envvariable.New),
		wire.NewSet(mount.New),
		wire.NewSet(worker.NewRegistry),
		wire.NewSet(http.New),
		wire.NewSet(logger.New),
		wire.NewSet(ws.New),
		wire.NewSet(scheduler.New),
		wire.NewSet(stats.New),
		wire.NewSet(newApp, newConfig),
	)))
}
