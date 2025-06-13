package router

import (
	"flag"
	"os"
	"testing"

	"gateway/internal/config"
	"github.com/gin-gonic/gin"
	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
)

func TestProvideServiceConfig(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "cfg*.json")
	if err != nil {
		t.Fatal(err)
	}
	data := `{"version":3,"name":"test","port":1234,"host":["http://localhost"],"endpoints":[]}`
	if _, err := tmp.WriteString(data); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	oldArgs := os.Args
	oldFS := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-c", tmp.Name()}

	cfg := &config.Config{}
	svcCfg, err := ProvideServiceConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if svcCfg.Port != 1234 {
		t.Errorf("expected 1234 got %d", svcCfg.Port)
	}

	flag.CommandLine = oldFS
	os.Args = oldArgs
}

func TestProvideServiceConfigOverride(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "cfg*.json")
	if err != nil {
		t.Fatal(err)
	}
	data := `{"version":3,"name":"test","port":1234,"host":["http://localhost"],"endpoints":[]}`
	if _, err := tmp.WriteString(data); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	oldArgs := os.Args
	oldFS := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-c", tmp.Name(), "-p", "9999"}

	cfg := &config.Config{}
	svcCfg, err := ProvideServiceConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if svcCfg.Port != 9999 {
		t.Errorf("expected 9999 got %d", svcCfg.Port)
	}

	flag.CommandLine = oldFS
	os.Args = oldArgs
}

func TestProvideGinRouter(t *testing.T) {
	logger := logging.NoOp
	svcCfg := &luraConfig.ServiceConfig{Port: 1234}
	engine := ProvideGinRouter(logger, svcCfg, &config.Config{})
	if len(engine.Handlers) != 3 {
		t.Errorf("expected default handlers 3 got %d", len(engine.Handlers))
	}

	cfg := &config.Config{Sentry: config.SentryConfig{Enable: true}, Cluster: "c"}
	engine = ProvideGinRouter(logger, svcCfg, cfg)
	if len(engine.Handlers) != 5 {
		t.Errorf("expected 5 handlers with sentry, got %d", len(engine.Handlers))
	}
}

func TestProvideRouter(t *testing.T) {
	logger := logging.NoOp
	pf := proxyFactoryMock{}
	eng := gin.New()
	cfg := &config.Config{}
	r := ProvideRouter(logger, pf, cfg, eng)
	if r == nil {
		t.Fatal("nil router")
	}
}

type proxyFactoryMock struct{}

func (proxyFactoryMock) New(*luraConfig.EndpointConfig) (proxy.Proxy, error) {
	return proxy.NoopProxy, nil
}
