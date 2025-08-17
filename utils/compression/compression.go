package compression

import (
	"io"
	"strings"
)

// DefaultCompressor the default compressors.
var DefaultCompressor = NewCompressors()

// CompressionFormat represents a compression format enumeration.
type CompressionFormat string

// Compressor abstracts the interface for a compression handler.
type Compressor interface {
	Compress(w io.Writer, src io.Reader) (int64, error)
	Decompress(reader io.ReadCloser) (io.ReadCloser, error)
}

// Compressors is a general helper for web compression.
type Compressors struct {
	acceptEncoding string
	compressors    map[CompressionFormat]Compressor
}

// NewCompressors create a Compressors instance.
func NewCompressors() *Compressors {
	compressors := map[CompressionFormat]Compressor{
		EncodingGzip:    GzipCompressor{},
		EncodingDeflate: DeflateCompressor{},
		EncodingZstd:    ZstdCompressor{},
	}

	return &Compressors{
		acceptEncoding: strings.Join([]string{
			string(EncodingZstd),
			string(EncodingGzip),
			string(EncodingDeflate),
		}, ", "),
		compressors: compressors,
	}
}

// AcceptEncoding returns the Accept-Encoding header with supported compression encodings.
func (c Compressors) AcceptEncoding() string {
	return c.acceptEncoding
}

// IsEncodingSupported checks if the input encoding is supported.
func (c Compressors) IsEncodingSupported(encoding string) bool {
	return c.FindSupportedEncoding(encoding) != ""
}

// FindSupportedEncoding returns the supported encoding from the input string.
func (c Compressors) FindSupportedEncoding(encoding string) CompressionFormat {
	_, ok := c.compressors[CompressionFormat(encoding)]
	if ok {
		return CompressionFormat(encoding)
	}

	for enc := range strings.SplitSeq(encoding, ",") {
		enc = strings.ToLower(strings.TrimSpace(enc))

		_, ok := c.compressors[CompressionFormat(enc)]
		if ok {
			return CompressionFormat(enc)
		}
	}

	return ""
}

// Compress writes compressed data.
func (c Compressors) Compress(w io.Writer, encoding string, data io.Reader) (int64, error) {
	format := c.FindSupportedEncoding(encoding)

	if format != "" {
		compressor, ok := c.compressors[format]
		if ok {
			return compressor.Compress(w, data)
		}
	}

	return io.Copy(w, data)
}

// Decompress reads and decompresses the reader with equivalent the content encoding.
func (c Compressors) Decompress(reader io.ReadCloser, encoding string) (io.ReadCloser, error) {
	format := c.FindSupportedEncoding(encoding)

	if format != "" {
		compressor, ok := c.compressors[format]
		if ok {
			return compressor.Decompress(reader)
		}
	}

	return reader, nil
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
