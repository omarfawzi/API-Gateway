package servers

import (
	"context"
	"net/http"
	"testing"

	"github.com/luraproject/lura/v2/transport/http/server/plugin"
)

type fakeRegisterer struct{ called int }

func (f *fakeRegisterer) RegisterHandlers(func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f.called++
}

func TestHandlerRegistry_RegisterHandlers(t *testing.T) {
	r1 := &fakeRegisterer{}
	r2 := &fakeRegisterer{}

	reg := NewHandlerRegistry([]plugin.Registerer{r1, r2})
	reg.RegisterHandlers()

	if r1.called != 1 || r2.called != 1 {
		t.Errorf("expected each registerer to be called once, got %d and %d", r1.called, r2.called)
	}
}
