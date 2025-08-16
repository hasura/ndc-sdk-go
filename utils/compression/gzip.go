package compression

import (
	"io"

	"github.com/klauspost/compress/gzip"
)

const (
	EncodingGzip = "gzip"
)

// GzipCompressor implements the compression handler for gzip encoding.
type GzipCompressor struct{}

// Compress the reader content with gzip encoding.
func (gc GzipCompressor) Compress(w io.Writer, src io.Reader) (int64, error) {
	zw := gzip.NewWriter(w)

	size, err := io.Copy(zw, src)
	_ = zw.Close()

	return size, err
}

// Decompress the reader content with gzip encoding.
func (gc GzipCompressor) Decompress(reader io.ReadCloser) (io.ReadCloser, error) {
	compressionReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return readCloserWrapper{
		CompressionReader: compressionReader,
		OriginalReader:    reader,
	}, nil
}
