package main

import (
	"gateway/internal"
	"gateway/internal/config"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"

	sentry "github.com/getsentry/sentry-go"
)

const (
	exitCodeOK = iota
	exitCodeError
	sentryFlushMaxWait = 5 * time.Second
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return exitCodeError
	}

	if cfg.Sentry.Enable {
		if err = sentry.Init(sentry.ClientOptions{
			Dsn:         cfg.Sentry.Dsn,
			Environment: cfg.Environment,
			Debug:       cfg.Sentry.Debug,
		}); err != nil {
			log.Printf("Sentry init failed: %v", err)
		} else {
			defer sentry.Flush(sentryFlushMaxWait)
		}
	}

	luraConfig.RoutingPattern = luraConfig.ColonRouterPatternBuilder

	svcCfg, logger, engine, err := buildCoreServices(cfg)
	if err != nil {
		return exitCodeError
	}

	err = registerHandlers(cfg)
	if err != nil {
		log.Printf("failed to register handlers: %v", err)
		return exitCodeError
	}
	routerFactory, err := internal.SetupRouter(cfg, engine, logger)
	if err != nil {
		log.Printf("failed to build router: %v", err)
		return exitCodeError
	}

	log.Printf("🚀 Gateway running on port %d", svcCfg.Port)
	routerFactory.Run(*svcCfg)

	return exitCodeOK
}

func buildCoreServices(cfg *config.Config) (*luraConfig.ServiceConfig, logging.Logger, *gin.Engine, error) {
	svcCfg, err := internal.ProvideServiceConfig(cfg)
	if err != nil {
		log.Printf("failed to build service config: %v", err)
		return nil, nil, nil, err
	}

	logger, err := logging.NewLogger(cfg.Logger.Level, os.Stdout, "[GATEWAY]")
	if err != nil {
		log.Printf("failed to build logger: %v", err)
		return nil, nil, nil, err
	}

	engine, err := internal.ProvideEngine(logger, svcCfg, cfg)
	if err != nil {
		log.Printf("failed to build engine: %v", err)
		return nil, nil, nil, err
	}

	return svcCfg, logger, engine, nil
}

func registerHandlers(cfg *config.Config) error {
	serverHandlerRegistry, err := internal.ProvideServerHandlerRegistry(cfg)
	if err != nil {
		return err
	}
	serverHandlerRegistry.RegisterHandlers()

	clientHandlerRegistry, err := internal.ProvideClientHandlerRegistry()
	if err != nil {
		return err
	}
	clientHandlerRegistry.RegisterHandlers()

	proxyHandlerRegistry, err := internal.ProvideProxyHandlerRegistry(cfg)
	if err != nil {
		return err
	}
	proxyHandlerRegistry.RegisterHandlers()

	return nil
}
