package clients

import (
	"net/http"
	"testing"
)

func TestNewHTTPClient_HTTP1(t *testing.T) {
	c := newHTTPClient("1")
	tr, ok := c.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected transport")
	}
	if tr.ForceAttemptHTTP2 {
		t.Errorf("expected http2 disabled")
	}
}

func TestNewHTTPClient_HTTP2(t *testing.T) {
	c := newHTTPClient("2")
	tr, ok := c.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected transport")
	}
	if !tr.ForceAttemptHTTP2 {
		t.Errorf("expected http2 enabled")
	}
}
