package connector

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/hasura/ndc-sdk-go/v2/utils/compression"
)

const (
	acceptEncodingHeader  = "Accept-Encoding"
	contentEncodingHeader = "Content-Encoding"
)

// define a custom response write to capture response information for logging.
type customResponseWriter struct {
	http.ResponseWriter

	contentEncoding string
	statusCode      int
	bodyLength      int
	body            []byte
}

var _ http.ResponseWriter = (*customResponseWriter)(nil)

func newCustomResponseWriter(r *http.Request, w http.ResponseWriter) *customResponseWriter {
	var contentEncoding string

	acceptEncoding := r.Header.Get(acceptEncodingHeader)

	if acceptEncoding != "" {
		for enc := range strings.SplitSeq(acceptEncoding, ",") {
			enc = strings.ToLower(strings.TrimSpace(enc))

			if compression.DefaultCompressor.IsEncodingSupported(enc) {
				contentEncoding = enc

				break
			}
		}
	}

	return &customResponseWriter{
		ResponseWriter:  w,
		contentEncoding: contentEncoding,
	}
}

func (cw *customResponseWriter) WriteHeader(statusCode int) {
	cw.statusCode = statusCode

	if cw.contentEncoding != "" {
		cw.ResponseWriter.Header().Set(contentEncodingHeader, cw.contentEncoding)
	}

	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *customResponseWriter) Write(body []byte) (int, error) {
	cw.body = body
	cw.bodyLength = len(body)

	written, err := compression.DefaultCompressor.Compress(
		cw.ResponseWriter,
		cw.contentEncoding,
		bytes.NewReader(body),
	)

	return int(written), err
}
