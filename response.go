package httpadapt

import (
	"bytes"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// ResponseWriter buffers response bytes
type ResponseWriter struct{ buf *bytes.Buffer }

// NewResponseWriter creates a writer that implements http.ResponseWriter that buffers the
// response body to be send as a response to the lambda proxy event
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{buf: bytes.NewBuffer(nil)}
}

// WriteHeader sends an HTTP response header with the provided status code.
func (w *ResponseWriter) WriteHeader(statusCode int) {}

// Header returns the header map that will be sent by WriteHeader.
func (w *ResponseWriter) Header() (h http.Header) { return }

// Write writes the data to the connection as part of an HTTP reply.
func (w *ResponseWriter) Write(b []byte) (int, error) { return w.buf.Write(b) }

// ProxyResponse returns the resulting response to handling the API Gateway Proxy event
func (w *ResponseWriter) ProxyResponse() (out events.APIGatewayProxyResponse, err error) {
	return
}
