//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/kiritoxkiriko/comical-tool/internal/server"
	"github.com/kiritoxkiriko/comical-tool/pkg/app"
	"github.com/kiritoxkiriko/comical-tool/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var serverSet = wire.NewSet(
	server.NewTask,
)

// build App
func newApp(
	task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(task),
		app.WithName("demo-task"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		serverSet,
		newApp,
	))
}
