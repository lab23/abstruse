//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/lab23/abstruse/worker/app"
	"github.com/lab23/abstruse/worker/logger"
)

func CreateApp() (*application, error) {
	panic(wire.Build(wire.NewSet(
		wire.NewSet(logger.New),
		wire.NewSet(app.NewApp),
		wire.NewSet(newApplication, newConfig),
	)))
}
