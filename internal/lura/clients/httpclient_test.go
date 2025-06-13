package clients

import (
	luraConfig "github.com/luraproject/lura/v2/config"
	"testing"
)

func TestParseHTTPVersion(t *testing.T) {
	be := &luraConfig.Backend{ExtraConfig: map[string]interface{}{
		httpVersionNamespace: map[string]interface{}{"version": "1.1"},
	}}
	if v := parseHTTPVersion(be); v != "1.1" {
		t.Errorf("expected 1.1 got %s", v)
	}

	be = &luraConfig.Backend{}
	if v := parseHTTPVersion(be); v != "2" {
		t.Errorf("expected default 2 got %s", v)
	}
}
