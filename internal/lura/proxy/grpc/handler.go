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
	"net/url"
	"strings"
)

type grpcHandler struct {
	conn          *grpc.ClientConn
	fullMethod    string
	inputFactory  func() proto.Message
	outputFactory func() proto.Message
}

func (h grpcHandler) asProxy(cfg *config.Backend) proxy.Proxy {
	return func(ctx context.Context, r *proxy.Request) (*proxy.Response, error) {
		body, err := readRequestBody(r)
		if err != nil {
			return nil, err
		}

		headerParams := extractHeaders(r.Headers, cfg.HeadersToPass)
		if len(headerParams) > 0 {
			ctx = metadata.NewOutgoingContext(ctx, metadata.New(headerParams))
		}

		params := mergeParams(r.Params, extractQueryParams(r.Query, cfg.QueryStringsToPass))

		input := h.inputFactory()
		if err := decodeRequestBody(body, input); err != nil {
			return nil, err
		}

		if err := populateProtoMessageFromParams(input, params); err != nil {
			return nil, err
		}

		output := h.outputFactory()
		if err := h.conn.Invoke(ctx, h.fullMethod, input, output); err != nil {
			return nil, err
		}

		return buildProxyResponse(output)
	}
}

func readRequestBody(r *proxy.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

func extractHeaders(headers map[string][]string, keys []string) map[string]string {
	result := make(map[string]string)
	for _, key := range keys {
		if values, ok := headers[key]; ok && len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

func extractQueryParams(query url.Values, keys []string) map[string][]string {
	result := make(map[string][]string)
	for _, key := range keys {
		if values, ok := query[key]; ok {
			result[strings.ToLower(key)] = values
		}
	}
	return result
}

func decodeRequestBody(body []byte, input proto.Message) error {
	if len(body) == 0 {
		return nil
	}
	return protojson.Unmarshal(body, input)
}

func buildProxyResponse(output proto.Message) (*proxy.Response, error) {
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

func mergeParams(pathParams map[string]string, queryParams map[string][]string) map[string][]string {
	merged := make(map[string][]string, len(pathParams)+len(queryParams))
	for k, v := range pathParams {
		lk := strings.ToLower(k)
		merged[lk] = []string{v}
	}
	for k, v := range queryParams {
		lk := strings.ToLower(k)
		existing := merged[lk]
		merged[lk] = append(existing, v...)
	}
	return merged
}
