package clients

import (
	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/transport/http/client"
	"github.com/luraproject/lura/v2/transport/http/client/plugin"
)

func ProvideHTTPRequestExecutor(
	logger logging.Logger,
) func(*luraConfig.Backend) client.HTTPRequestExecutor {
	return plugin.HTTPRequestExecutor(logger, func(*luraConfig.Backend) client.HTTPRequestExecutor {
		return client.DefaultHTTPRequestExecutor(client.NewHTTPClient)
	})
}

func ProvideHandlers() []plugin.Registerer {
	return []plugin.Registerer{grpcRegisterer{}}
}
