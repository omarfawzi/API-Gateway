package errors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareErrorResponse(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "value")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("bad"))
		if err != nil {
			return
		}
	})

	handler, err := NewServerHandler().middleware(context.Background(), nil, next)
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d got %d", http.StatusInternalServerError, rr.Code)
	}
	expected := "{\"errors\":{\"message\":\"Internal Server Error\"}}\n"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
	if rr.Header().Get("X-Test") != "" {
		t.Errorf("unexpected header forwarded")
	}
}

func TestMiddlewareSuccess(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("hello"))
		if err != nil {
			return
		}
	})

	handler, err := NewServerHandler().middleware(context.Background(), nil, next)
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d got %d", http.StatusOK, rr.Code)
	}
	if rr.Header().Get("X-Test") != "ok" {
		t.Errorf("header not forwarded")
	}
	if rr.Body.String() != "hello" {
		t.Errorf("body not forwarded: %s", rr.Body.String())
	}
}
