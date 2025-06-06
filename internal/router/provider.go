package router

import (
	"context"
	"flag"
	"gateway/internal/config"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	sentryGin "github.com/getsentry/sentry-go/gin"
	sentryhttp "github.com/getsentry/sentry-go/http"
	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"github.com/luraproject/lura/v2/router"
	krakendGin "github.com/luraproject/lura/v2/router/gin"
	"github.com/luraproject/lura/v2/transport/http/server"
	handlerPlugins "github.com/luraproject/lura/v2/transport/http/server/plugin"
)

const (
	configFile = "config/config.json"
)

func ProvideRouter(
	logger logging.Logger,
	proxyFactory proxy.Factory,
	cfg *config.Config,
	engine *gin.Engine,
) router.Router {
	runServer := handlerPlugins.New(
		logger,
		func(
			ctx context.Context,
			luraCfg luraConfig.ServiceConfig,
			handler http.Handler,
		) error {

			if !cfg.Sentry.Enable {
				return server.RunServer(ctx, luraCfg, handler)
			}

			sentryHandler := sentryhttp.New(sentryhttp.Options{
				Repanic: true,
			})

			wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hub := sentry.CurrentHub().Clone()
				hub.Scope().SetTag("cluster", cfg.Cluster)

				ctxWithHub := sentry.SetHubOnContext(r.Context(), hub)
				r = r.WithContext(ctxWithHub)

				sentryHandler.Handle(handler).ServeHTTP(w, r)
			})

			return server.RunServer(ctx, luraCfg, wrappedHandler)
		})

	return krakendGin.NewFactory(
		krakendGin.Config{
			Engine:         engine,
			HandlerFactory: krakendGin.CustomErrorEndpointHandler(logger, server.DefaultToHTTPError),
			ProxyFactory:   proxyFactory,
			Logger:         logger,
			RunServer:      krakendGin.RunServerFunc(runServer),
		},
	).New()
}

func ProvideServiceConfig(cfg *config.Config) (*luraConfig.ServiceConfig, error) {
	port := flag.Int("p", cfg.Port, "Port of the app")
	luraConfigFile := flag.String("c", configFile, "Lura configuration file. Default: config/config.json")
	flag.Parse()

	parser := luraConfig.NewParser()
	conf, err := parser.Parse(*luraConfigFile)
	if err != nil {
		log.Printf("failed to parse Lura config: %v", err)
		return nil, err
	}

	if port != nil && *port != 0 {
		conf.Port = *port
	}

	return &conf, nil
}

func ProvideGinRouter(logger logging.Logger, svcCfg *luraConfig.ServiceConfig, cfg *config.Config) *gin.Engine {
	engine := krakendGin.NewEngine(*svcCfg, krakendGin.EngineOptions{
		Logger: logger,
	})

	if cfg.Sentry.Enable {
		engine.Use(sentryGin.New(sentryGin.Options{
			Repanic: true,
		}))

		engine.Use(func(ctx *gin.Context) {
			if hub := sentryGin.GetHubFromContext(ctx); hub != nil {
				hub.Scope().SetTag("cluster", cfg.Cluster)
			}
			ctx.Next()
		})
	}

	return engine
}
