//nolint:staticcheck // use deprecated protoreflect packages for reflection
package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	luraConfig "github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"google.golang.org/grpc"
	reflpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

const grpcNamespace = "github.com/devopsfaith/krakend/grpc"

type grpcConfig struct {
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

// GRPCBackendFactory wraps a backend factory to provide gRPC support using
// reflection when the backend contains a grpcNamespace configuration.
func GRPCBackendFactory(logger logging.Logger, next proxy.BackendFactory) proxy.BackendFactory {
	return func(remote *luraConfig.Backend) proxy.Proxy {
		cfgRaw, ok := remote.ExtraConfig[grpcNamespace]
		if !ok {
			return next(remote)
		}
		data, err := json.Marshal(cfgRaw)
		if err != nil {
			logger.Error("[gRPC] marshal config:", err.Error())
			return next(remote)
		}
		var cfg grpcConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			logger.Error("[gRPC] parse config:", err.Error())
			return next(remote)
		}
		if cfg.Endpoint == "" || cfg.Method == "" {
			logger.Error("[gRPC] missing endpoint or method")
			return next(remote)
		}
		service, method, err := splitMethod(cfg.Method)
		if err != nil {
			logger.Error("[gRPC] invalid method:", err.Error())
			return next(remote)
		}
		conn, err := grpc.Dial(cfg.Endpoint, grpc.WithInsecure())
		if err != nil {
			logger.Error("[gRPC] dial:", err.Error())
			return next(remote)
		}
		refClient := grpcreflect.NewClient(context.Background(), reflpb.NewServerReflectionClient(conn))
		svcDesc, err := refClient.ResolveService(service)
		if err != nil {
			logger.Error("[gRPC] resolve service:", err.Error())
			return next(remote)
		}
		mDesc := svcDesc.FindMethodByName(method)
		if mDesc == nil {
			logger.Error("[gRPC] method not found:", cfg.Method)
			return next(remote)
		}
		stub := grpcdynamic.NewStub(conn)
		return func(ctx context.Context, r *proxy.Request) (*proxy.Response, error) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			msg := dynamic.NewMessage(mDesc.GetInputType())
			if len(body) > 0 {
				if err := msg.UnmarshalJSON(body); err != nil {
					return nil, err
				}
			}
			respMsg, err := stub.InvokeRpc(ctx, mDesc, msg)
			if err != nil {
				return nil, err
			}
			out, err := respMsg.(interface{ MarshalJSON() ([]byte, error) }).MarshalJSON()
			if err != nil {
				return nil, err
			}
			var data map[string]interface{}
			if err := json.Unmarshal(out, &data); err != nil {
				return nil, err
			}
			return &proxy.Response{
				Data:       data,
				IsComplete: true,
				Metadata: proxy.Metadata{
					Headers: map[string][]string{
						"Content-Type": {"application/json"},
					},
				},
			}, nil
		}
	}
}

func splitMethod(full string) (string, string, error) {
	parts := strings.Split(full, "/")
	if len(parts) != 2 {
		return "", "", errors.New("method must be service/method")
	}
	return parts[0], parts[1], nil
}
