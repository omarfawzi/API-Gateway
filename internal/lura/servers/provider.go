package servers

import (
	"gateway/internal/errors"

	"github.com/luraproject/lura/v2/transport/http/server/plugin"
)

func ProvideHandlers(errorsHandler *errors.ServerHandler) []plugin.Registerer {
	return []plugin.Registerer{errorsHandler}
}
