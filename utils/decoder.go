package utils

import (
	"bytes"
	"github.com/klauspost/compress/gzip"
	"io"
)

func DecodeGzipBody(body []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
