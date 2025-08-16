package compression

import (
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

const (
	EncodingZstd = "zstd"
)

// ZstdCompressor implements the compression handler for zstandard encoding.
type ZstdCompressor struct{}

// Compress the reader content with gzip encoding.
func (zc ZstdCompressor) Compress(w io.Writer, src io.Reader) (int64, error) {
	zw, err := zstd.NewWriter(w)
	if err != nil {
		return 0, fmt.Errorf("failed to create the zstd writer: %w", err)
	}

	size, err := io.Copy(zw, src)
	_ = zw.Close()

	return size, err
}

// Decompress the reader content with gzip encoding.
func (zc ZstdCompressor) Decompress(reader io.ReadCloser) (io.ReadCloser, error) {
	compressionReader, err := zstd.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return readCloserWrapper{
		CompressionReader: &zstdReader{
			Decoder: compressionReader,
		},
		OriginalReader: reader,
	}, nil
}

type zstdReader struct {
	*zstd.Decoder
}

func (zd *zstdReader) Close() error {
	zd.Decoder.Close()

	return nil
}
