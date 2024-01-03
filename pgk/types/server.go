package types

import (
	"io"
	"net/http"
)

// Request wrapper
type RequestWrapper struct {
	Body   []byte      `cbor:"body"`
	Method string      `cbor:"method"`
	URL    string      `cbor:"url"`
	Header http.Header `cbor:"header"`
}

type ResponseWrapper struct {
	Body       []byte      `cbor:"body"`
	StatusCode int         `cbor:"statusCode"`
	Header     http.Header `cbor:"header"`
}

func NewRequestWrapperFromHttpRequest(httpReq *http.Request) (*RequestWrapper, error) {

	bodyBuffer, err := io.ReadAll(httpReq.Body)
	if err != nil {
		return nil, err
	}

	req := RequestWrapper{
		Body:   bodyBuffer,
		Method: httpReq.Method,
		URL:    httpReq.URL.String(),
		Header: httpReq.Header,
	}

	return &req, nil
}
