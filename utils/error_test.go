package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONError(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteJSONError(rr, http.StatusNotFound)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d got %d", http.StatusNotFound, rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("unexpected content type %s", ct)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Errors.Message != http.StatusText(http.StatusNotFound) {
		t.Errorf("unexpected message %s", resp.Errors.Message)
	}
}

func TestNewHTTPResponseError(t *testing.T) {
	errObj := NewHTTPResponseError(http.StatusBadRequest)
	if errObj.Code != http.StatusBadRequest {
		t.Errorf("expected code %d got %d", http.StatusBadRequest, errObj.Code)
	}
	if errObj.Encoding() != "application/json" {
		t.Errorf("unexpected encoding %s", errObj.Encoding())
	}
	if errObj.Msg == "" {
		t.Errorf("expected message not empty")
	}
}

func TestToJSONFallback(t *testing.T) {
	ch := make(chan int)
	out := toJSON(ch)
	if out != "{}" {
		t.Errorf("expected fallback '{}', got %s", out)
	}
}
