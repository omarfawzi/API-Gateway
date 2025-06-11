package grpc

import (
	"context"
	"encoding/json"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
)

type grpcHandler struct {
	conn          *grpc.ClientConn
	fullMethod    string
	inputFactory  func() proto.Message
	outputFactory func() proto.Message
}

func (h grpcHandler) asProxy(cfg *config.Backend) proxy.Proxy {
	return func(ctx context.Context, r *proxy.Request) (*proxy.Response, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		input := h.inputFactory()

		headerParams := make(map[string]string)
		for _, key := range cfg.HeadersToPass {
			if values, ok := r.Headers[key]; ok && len(values) > 0 {
				headerParams[key] = values[0]
			}
		}

		if len(headerParams) > 0 {
			ctx = metadata.NewOutgoingContext(ctx, metadata.New(headerParams))
		}

		queryParams := make(map[string]string)
		for _, key := range cfg.QueryStringsToPass {
			if r.Query.Has(key) {
				queryParams[key] = r.Query.Get(key)
			}
		}

		if len(body) > 0 {
			if err := protojson.Unmarshal(body, input); err != nil {
				return nil, err
			}
		}

		if err := populateProtoMessageFromParams(input, queryParams); err != nil {
			return nil, err
		}

		output := h.outputFactory()
		if err := h.conn.Invoke(ctx, h.fullMethod, input, output); err != nil {
			return nil, err
		}

		jsonBytes, err := protojson.Marshal(output)
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
