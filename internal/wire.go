//go:build wireinject
// +build wireinject

package internal

import (
	"gateway/internal/config"
	"gateway/internal/errors"
	"gateway/internal/lura/clients"
	"gateway/internal/lura/proxy"
	"gateway/internal/lura/servers"
	"gateway/internal/router"

	"github.com/gin-gonic/gin"

	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"

	"github.com/google/wire"
	luraRouter "github.com/luraproject/lura/v2/router"
)

func SetupRouter(cfg *config.Config, engine *gin.Engine, logger logging.Logger) (factory luraRouter.Router, lura error) {
	wire.Build(
		router.ProvideRouter,
		proxy.ProvideProxyFactory,
		clients.ProvideHTTPRequestExecutor,
	)
	return nil, nil
}

func ProvideServiceConfig(cfg *config.Config) (*luraConfig.ServiceConfig, error) {
	wire.Build(
		router.ProvideServiceConfig,
	)

	return nil, nil
}

func ProvideEngine(logger logging.Logger, svcCfg *luraConfig.ServiceConfig, cfg *config.Config) (*gin.Engine, error) {
	wire.Build(
		router.ProvideGinRouter,
	)

	return nil, nil
}

func ProvideServerHandlerRegistry(cfg *config.Config) (*servers.HandlerRegistry, error) {
	wire.Build(
		errors.ProvideServerHandler,
		servers.ProvideHandlers,
		servers.NewHandlerRegistry,
	)
	return nil, nil
}

func ProvideClientHandlerRegistry() (*clients.HandlerRegistry, error) {
	wire.Build(
		clients.ProvideHandlers,
		clients.NewHandlerRegistry,
	)
	return nil, nil
}

func ProvideProxyHandlerRegistry(cfg *config.Config) (*proxy.HandlerRegistry, error) {
	wire.Build(
		proxy.ProvideHandlers,
		proxy.NewHandlerRegistry,
	)
	return nil, nil
}
