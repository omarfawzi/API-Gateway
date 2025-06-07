package utils

import (
	"bytes"
	"github.com/klauspost/compress/gzip"
	"io"
)

func DecodeGzipBody(body []byte) (data []byte, err error) {
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := reader.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	data, err = io.ReadAll(reader)
	return
}
