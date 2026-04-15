package webview

import (
	"io"
	"net/http"
)

// ResponseWriter is the minimal response writer surface used by the application stubs.
type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

// Request mirrors the internal request interface used by Wails' asset server bridge.
type Request interface {
	URL() (string, error)
	Method() (string, error)
	Header() (http.Header, error)
	Body() (io.ReadCloser, error)
	Response() ResponseWriter
	Close() error
}
