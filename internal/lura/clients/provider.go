package clients

import (
	"context"
	"net/http"

	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/transport/http/client"
	"github.com/luraproject/lura/v2/transport/http/client/plugin"
)

func ProvideHTTPRequestExecutor(
	logger logging.Logger,
) func(*luraConfig.Backend) client.HTTPRequestExecutor {
	return plugin.HTTPRequestExecutor(logger, func(be *luraConfig.Backend) client.HTTPRequestExecutor {
		version := parseHTTPVersion(be)
		return client.DefaultHTTPRequestExecutor(func(context.Context) *http.Client {
			return newHTTPClient(version)
		})
	})
}

func ProvideHandlers() []plugin.Registerer {
	return []plugin.Registerer{}
}
