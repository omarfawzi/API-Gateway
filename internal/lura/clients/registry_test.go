package clients

import (
	"context"
	"net/http"
	"testing"

	"github.com/luraproject/lura/v2/transport/http/client/plugin"
)

type fakeRegisterer struct{ called int }

func (f *fakeRegisterer) RegisterClients(func(string, func(context.Context, map[string]interface{}) (http.Handler, error))) {
	f.called++
}

func TestHandlerRegistry_RegisterHandlers(t *testing.T) {
	r1 := &fakeRegisterer{}
	r2 := &fakeRegisterer{}

	registry := NewHandlerRegistry([]plugin.Registerer{r1, r2})
	registry.RegisterHandlers()

	if r1.called != 1 || r2.called != 1 {
		t.Errorf("expected each registerer to be called once, got %d and %d", r1.called, r2.called)
	}
}

func TestProvideHandlersEmpty(t *testing.T) {
	if len(ProvideHandlers()) != 0 {
		t.Errorf("expected empty handlers slice")
	}
}
