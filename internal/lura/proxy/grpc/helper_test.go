package grpc

import (
	"net/url"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestExtractHeaders(t *testing.T) {
	headers := map[string][]string{"A": {"1"}, "B": {"2"}}
	out := extractHeaders(headers, []string{"B"})
	if out["B"] != "2" {
		t.Errorf("expected B=2")
	}
}

func TestExtractQueryParams(t *testing.T) {
	q := url.Values{"A": {"1"}, "B": {"2"}}
	out := extractQueryParams(q, []string{"A"})
	if v := out["a"]; len(v) != 1 || v[0] != "1" {
		t.Errorf("unexpected %v", out)
	}
}

func TestDecodeRequestBody(t *testing.T) {
	msg := &structpb.Struct{}
	if err := decodeRequestBody([]byte(`{"foo":"bar"}`), msg); err != nil {
		t.Fatal(err)
	}
	if msg.Fields["foo"].GetStringValue() != "bar" {
		t.Errorf("decoded value wrong")
	}
	if err := decodeRequestBody(nil, msg); err != nil {
		t.Fatal(err)
	}
}

func TestBuildProxyResponse(t *testing.T) {
	msg := &structpb.Struct{Fields: map[string]*structpb.Value{"a": structpb.NewStringValue("b")}}
	resp, err := buildProxyResponse(msg)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["a"].(string) != "b" {
		t.Errorf("wrong data")
	}
	if resp.Metadata.Headers["Content-Type"][0] != "application/json" {
		t.Errorf("wrong header")
	}
}

func TestMergeParams(t *testing.T) {
	path := map[string]string{"ID": "1"}
	query := map[string][]string{"id": {"2"}}
	out := mergeParams(path, query)
	if len(out["id"]) != 2 {
		t.Errorf("expected merged slice")
	}
}

func TestParseServiceMethod(t *testing.T) {
	svc, m, err := parseServiceMethod("svc/Method")
	if err != nil || svc != "svc" || m != "Method" {
		t.Errorf("unexpected result %s %s %v", svc, m, err)
	}
	if _, _, err := parseServiceMethod("bad"); err == nil {
		t.Errorf("expected error")
	}
}
