package proxy

import (
	"gateway/internal/lura/proxy/grpc"
	cb "github.com/krakendio/krakend-circuitbreaker/v2/gobreaker/proxy"
	martian "github.com/krakendio/krakend-martian/v2"
	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	proxyPlugins "github.com/luraproject/lura/v2/proxy/plugin"
	"github.com/luraproject/lura/v2/transport/http/client"
	"github.com/luraproject/lura/v2/transport/http/client/plugin"
)

func ProvideProxyFactory(
	logger logging.Logger,
	httpRequestExecutor func(*luraConfig.Backend) client.HTTPRequestExecutor,
) proxy.Factory {
	executorFactory := plugin.HTTPRequestExecutor(logger, httpRequestExecutor)
	martianFactory := martian.NewConfiguredBackendFactory(logger, executorFactory)
	grpcFactory := grpc.NewBackendFactory(logger, martianFactory)
	backendFactory := cb.BackendFactory(grpcFactory, logger)

	return proxy.NewDefaultFactory(backendFactory, logger)
}

func ProvideHandlers() []proxyPlugins.Registerer {
	return []proxyPlugins.Registerer{}
}
