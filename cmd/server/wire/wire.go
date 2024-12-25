//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/kiritoxkiriko/comical-tool/internal/handler"
	"github.com/kiritoxkiriko/comical-tool/internal/repository"
	"github.com/kiritoxkiriko/comical-tool/internal/server"
	"github.com/kiritoxkiriko/comical-tool/internal/service"
	"github.com/kiritoxkiriko/comical-tool/pkg/app"
	"github.com/kiritoxkiriko/comical-tool/pkg/jwt"
	"github.com/kiritoxkiriko/comical-tool/pkg/log"
	"github.com/kiritoxkiriko/comical-tool/pkg/server/http"
	"github.com/kiritoxkiriko/comical-tool/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
)

// build App
func newApp(
	httpServer *http.Server,
	job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
