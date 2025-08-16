package compression

import (
	"io"

	"github.com/klauspost/compress/flate"
)

const (
	EncodingDeflate = "deflate"
)

// DeflateCompressor implements the compression handler for deflate encoding.
type DeflateCompressor struct{}

// Compress the reader content with deflate encoding.
func (dc DeflateCompressor) Compress(w io.Writer, src io.Reader) (int64, error) {
	fw, err := flate.NewWriter(w, flate.DefaultCompression)
	if err != nil {
		return 0, err
	}

	size, copyErr := io.Copy(fw, src)
	closeErr := fw.Close()

	if copyErr != nil {
		return 0, copyErr
	}

	if closeErr != nil {
		return 0, closeErr
	}

	return size, nil
}

// Decompress the reader content with deflate encoding.
func (dc DeflateCompressor) Decompress(reader io.ReadCloser) (io.ReadCloser, error) {
	compressionReader := flate.NewReader(reader)

	return readCloserWrapper{
		CompressionReader: compressionReader,
		OriginalReader:    reader,
	}, nil
}
