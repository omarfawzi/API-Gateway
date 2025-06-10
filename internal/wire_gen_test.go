package internal

import (
	"reflect"
	"testing"

	"gateway/internal/config"
)

func TestProvideClientHandlerRegistry(t *testing.T) {
	r, err := ProvideClientHandlerRegistry()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("nil registry returned")
	}
	handlers := reflect.ValueOf(r).Elem().FieldByName("handlers")
	if handlers.Len() != 1 {
		t.Errorf("expected one handler, got %d", handlers.Len())
	}
}

func TestProvideProxyHandlerRegistry(t *testing.T) {
	r, err := ProvideProxyHandlerRegistry(&config.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("nil registry returned")
	}
	handlers := reflect.ValueOf(r).Elem().FieldByName("handlers")
	if handlers.Len() != 0 {
		t.Errorf("expected no handlers, got %d", handlers.Len())
	}
}
