package clients

import (
	"github.com/luraproject/lura/v2/transport/http/client/plugin"
)

type HandlerRegistry struct {
	handlers []plugin.Registerer
}

func NewHandlerRegistry(
	handlers []plugin.Registerer,
) *HandlerRegistry {
	return &HandlerRegistry{
		handlers: handlers,
	}
}

func (r *HandlerRegistry) RegisterHandlers() {
	for _, h := range r.handlers {
		h.RegisterClients(plugin.RegisterClient)
	}
}
