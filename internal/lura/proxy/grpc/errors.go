package grpc

import "errors"

var (
	errInvalidParams       = errors.New("[gRPC] invalid params")
	errInvalidMethod       = errors.New("[gRPC] method not found")
	errInvalidService      = errors.New("[gRPC] service not found")
	errInvalidMethodFormat = errors.New("[gRPC] method must be service/method")
	errInvalidJSON         = errors.New("[gRPC] failed to marshal response to JSON")
	errMissingGRPCConfig   = errors.New("[gRPC] missing gRPC config")
	errMissingMethodConfig = errors.New("[gRPC] missing method in config")
	errMissingDescriptor   = errors.New("[gRPC] missing descriptor_file in config")
	errUnsupportedKind     = errors.New("[gRPC] unsupported kind")
)
