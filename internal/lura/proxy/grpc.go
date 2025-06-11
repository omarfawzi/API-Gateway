package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/reflect/protoregistry"
	"io"
	"os"
	"strings"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

const Namespace = "plugins/grpc"

var (
	errInvalidMethod       = errors.New("[gRPC] method not found")
	errInvalidService      = errors.New("[gRPC] service not found")
	errInvalidMethodFormat = errors.New("[gRPC] method must be service/method")
	errInvalidJSON         = errors.New("[gRPC] failed to marshal response to JSON")
	errMissingGRPCConfig   = errors.New("[gRPC] missing gRPC config")
	errMissingMethodConfig = errors.New("[gRPC] missing method in config")
	errMissingDescriptor   = errors.New("[gRPC] missing descriptor_file in config")
)

type grpcConfig struct {
	Method         string `json:"method"`
	DescriptorFile string `json:"descriptor_file"`
}

func NewGrpcBackendFactory(logger logging.Logger, next proxy.BackendFactory) proxy.BackendFactory {
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

		return buildProxyHandler(conn, desc)
	}
}

func extractConfig(remote *config.Backend, logger logging.Logger) (*grpcConfig, error) {
	rawCfg, ok := remote.ExtraConfig[Namespace]
	if !ok {
		return nil, errMissingGRPCConfig
	}

	data, err := json.Marshal(rawCfg)
	if err != nil {
		logger.Error("[gRPC] marshal config:", err.Error())
		return nil, err
	}

	var cfg grpcConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		logger.Error("[gRPC] parse config:", err.Error())
		return nil, err
	}

	if cfg.Method == "" {
		logger.Error("[gRPC] missing method")
		return nil, errMissingMethodConfig
	}
	if cfg.DescriptorFile == "" {
		logger.Error("[gRPC] missing descriptor_file")
		return nil, errMissingDescriptor
	}

	return &cfg, nil
}

func methodDescriptorFromFile(
	descriptorPath, serviceName, methodName string,
) (protoreflect.MethodDescriptor, error) {
	fds, err := readDescriptorSet(descriptorPath)
	if err != nil {
		return nil, err
	}

	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return nil, fmt.Errorf("[gRPC] create file registry: %w", err)
	}

	svcDesc, err := findServiceDescriptor(files, fds, serviceName)
	if err != nil {
		return nil, err
	}

	methodDesc := svcDesc.Methods().ByName(protoreflect.Name(methodName))
	if methodDesc == nil {
		return nil, errInvalidMethod
	}

	return methodDesc, nil
}

func readDescriptorSet(path string) (*descriptorpb.FileDescriptorSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[gRPC] read descriptor file: %w", err)
	}

	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &fds); err != nil {
		return nil, fmt.Errorf("[gRPC] unmarshal descriptor: %w", err)
	}
	return &fds, nil
}

func findServiceDescriptor(
	files *protoregistry.Files,
	fds *descriptorpb.FileDescriptorSet,
	targetService string,
) (protoreflect.ServiceDescriptor, error) {
	for _, file := range fds.File {
		fd, err := files.FindFileByPath(file.GetName())
		if err != nil {
			continue
		}
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			if string(svc.FullName()) == targetService {
				return svc, nil
			}
		}
	}
	return nil, errInvalidService
}

func buildProxyHandler(conn *grpc.ClientConn, methodDesc protoreflect.MethodDescriptor) proxy.Proxy {
	return func(ctx context.Context, r *proxy.Request) (*proxy.Response, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		inputMsg := dynamicpb.NewMessage(methodDesc.Input())
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, inputMsg); err != nil {
				return nil, err
			}
		}

		outputMsg := dynamicpb.NewMessage(methodDesc.Output())

		fullMethod := fmt.Sprintf("/%s/%s", methodDesc.Parent().FullName(), methodDesc.Name())
		if err := conn.Invoke(ctx, fullMethod, inputMsg, outputMsg); err != nil {
			return nil, err
		}

		jsonBytes, err := protojson.Marshal(outputMsg)
		if err != nil {
			return nil, errInvalidJSON
		}

		var data map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &data); err != nil {
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

func parseServiceMethod(full string) (string, string, error) {
	parts := strings.Split(full, "/")
	if len(parts) != 2 {
		return "", "", errInvalidMethodFormat
	}
	return parts[0], parts[1], nil
}
