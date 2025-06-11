package grpc

import (
	"encoding/json"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
)

type grpcConfig struct {
	Method         string `json:"method"`
	DescriptorFile string `json:"descriptor_file"`
}

const Namespace = "plugins/grpc"

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
