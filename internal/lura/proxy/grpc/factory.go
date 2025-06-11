package grpc

import (
	"fmt"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
	"strings"
)

func NewBackendFactory(logger logging.Logger, next proxy.BackendFactory) proxy.BackendFactory {
	return func(remote *config.Backend) proxy.Proxy {
		cfg, err := extractConfig(remote, logger)
		if err != nil {
			return next(remote)
		}

		service, method, err := parseServiceMethod(cfg.Method)
		if err != nil {
			logger.Error("[gRPC] invalid method:", err.Error())
			return next(remote)
		}

		desc, err := methodDescriptorFromFile(cfg.DescriptorFile, service, method)
		if err != nil {
			logger.Error("[gRPC] resolve method:", err.Error())
			return next(remote)
		}

		addr := strings.TrimPrefix(remote.Host[0], "http://")
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error("[gRPC] dial:", err.Error())
			return next(remote)
		}

		fullMethod := fmt.Sprintf("/%s/%s", desc.Parent().FullName(), desc.Name())
		inputPrototype := dynamicpb.NewMessage(desc.Input())
		outputPrototype := dynamicpb.NewMessage(desc.Output())

		handler := grpcHandler{
			conn:       conn,
			fullMethod: fullMethod,
			inputFactory: func() proto.Message {
				return inputPrototype.ProtoReflect().New().Interface()
			},
			outputFactory: func() proto.Message {
				return outputPrototype.ProtoReflect().New().Interface()
			},
		}

		return handler.asProxy(remote)
	}
}
