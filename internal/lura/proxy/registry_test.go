package proxy

import (
	"testing"

	"github.com/luraproject/lura/v2/proxy/plugin"
)

type fakeRegisterer struct{ called int }

func (f *fakeRegisterer) RegisterModifiers(func(string, func(map[string]interface{}) func(interface{}) (interface{}, error), bool, bool)) {
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

func TestProvideHandlersEmpty(t *testing.T) {
	if len(ProvideHandlers()) != 0 {
		t.Errorf("expected empty handlers slice")
	}
}
