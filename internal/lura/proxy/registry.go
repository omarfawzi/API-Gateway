package proxy

import (
	"github.com/luraproject/lura/v2/proxy/plugin"
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
		h.RegisterModifiers(plugin.RegisterModifier)
	}
}
