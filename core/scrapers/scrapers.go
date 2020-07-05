package scrapers

import (
	"bytes"
	libgzip "compress/gzip"
)

func mustGzip(data []byte) []byte {
	gzipped, err := gzip(data)
	if err != nil {
		panic(err)
	}
	return gzipped
}

func gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := libgzip.NewWriter(&buf)

	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
