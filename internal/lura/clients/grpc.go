//nolint:staticcheck // use of deprecated protoreflect APIs is required for reflection support
package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/luraproject/lura/v2/transport/http/client/plugin"
	"google.golang.org/grpc"
	reflpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type grpcRegisterer struct{}

type grpcConfig struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

func (grpcRegisterer) RegisterClients(f func(string, func(context.Context, map[string]interface{}) (http.Handler, error))) {
	f("grpc", grpcHandler)
}

func grpcHandler(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	data, err := json.Marshal(extra)
	if err != nil {
		return nil, err
	}
	var cfg grpcConfig
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Endpoint == "" || cfg.Method == "" {
		return nil, errors.New("missing endpoint or method")
	}
	service, method, err := splitMethod(cfg.Method)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.DialContext(ctx, cfg.Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	refClient := grpcreflect.NewClient(ctx, reflpb.NewServerReflectionClient(conn))
	svcDesc, err := refClient.ResolveService(service)
	if err != nil {
		return nil, err
	}
	mDesc := svcDesc.FindMethodByName(method)
	if mDesc == nil {
		return nil, fmt.Errorf("method %s not found", method)
	}
	stub := grpcdynamic.NewStub(conn)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg := dynamic.NewMessage(mDesc.GetInputType())
		if len(body) > 0 {
			if err := msg.UnmarshalJSON(body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		respMsg, err := stub.InvokeRpc(r.Context(), mDesc, msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		out, err := respMsg.(interface{ MarshalJSON() ([]byte, error) }).MarshalJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
	}), nil
}

func splitMethod(full string) (string, string, error) {
	parts := strings.Split(full, "/")
	if len(parts) != 2 {
		return "", "", errors.New("method must be service/method")
	}
	return parts[0], parts[1], nil
}

var _ plugin.Registerer = (*grpcRegisterer)(nil)
