package clients

import (
	"crypto/tls"
	"net/http"
	"strings"

	luraConfig "github.com/luraproject/lura/v2/config"
)

const httpVersionNamespace = "backend/http/client"

func parseHTTPVersion(be *luraConfig.Backend) string {
	raw, ok := be.ExtraConfig[httpVersionNamespace]
	if !ok {
		return "2"
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return "2"
	}
	if v, ok := m["version"].(string); ok {
		return v
	}
	return "2"
}

func newHTTPClient(version string) *http.Client {
	tr := &http.Transport{ForceAttemptHTTP2: true}
	if strings.HasPrefix(version, "1") {
		tr.ForceAttemptHTTP2 = false
		tr.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	}
	return &http.Client{Transport: tr}
}
