package compression

import (
	"io"
	"strings"

	"github.com/hasura/ndc-sdk-go/utils"
)

// DefaultCompressor the default compressors.
var DefaultCompressor = NewCompressors()

// Compressor abstracts the interface for a compression handler.
type Compressor interface {
	Compress(w io.Writer, src io.Reader) (int64, error)
	Decompress(reader io.ReadCloser) (io.ReadCloser, error)
}

// Compressors is a general helper for web compression.
type Compressors struct {
	acceptEncoding string
	compressors    map[string]Compressor
}

// NewCompressors create a Compressors instance.
func NewCompressors() *Compressors {
	compressors := map[string]Compressor{
		EncodingGzip:    GzipCompressor{},
		EncodingDeflate: DeflateCompressor{},
		EncodingZstd:    ZstdCompressor{},
	}

	return &Compressors{
		acceptEncoding: strings.Join(utils.GetSortedKeys(compressors), ", "),
		compressors:    compressors,
	}
}

// AcceptEncoding returns the Accept-Encoding header with supported compression encodings.
func (c Compressors) AcceptEncoding() string {
	return c.acceptEncoding
}

// IsEncodingSupported checks if the input encoding is supported.
func (c Compressors) IsEncodingSupported(encoding string) bool {
	_, ok := c.compressors[encoding]

	return ok
}

// Compress writes compressed data.
func (c Compressors) Compress(w io.Writer, encoding string, data io.Reader) (int64, error) {
	compressor, ok := c.compressors[strings.ToLower(strings.TrimSpace(encoding))]
	if !ok {
		return io.Copy(w, data)
	}

	return compressor.Compress(w, data)
}

// Decompress reads and decompresses the reader with equivalent the content encoding.
func (c Compressors) Decompress(reader io.ReadCloser, encoding string) (io.ReadCloser, error) {
	compressor, ok := c.compressors[strings.ToLower(strings.TrimSpace(encoding))]
	if !ok {
		return reader, nil
	}

	return compressor.Decompress(reader)
}

type readCloserWrapper struct {
	CompressionReader io.ReadCloser
	OriginalReader    io.ReadCloser
}

func (rcw readCloserWrapper) Close() error {
	_ = rcw.OriginalReader.Close()

	return rcw.CompressionReader.Close()
}

func (rcw readCloserWrapper) Read(p []byte) (int, error) {
	return rcw.CompressionReader.Read(p)
}
