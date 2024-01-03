package sdk

import (
	"bytes"
	"net/http"
)

// Estructura privada para parsear el bloque de memoria desde el servidor
type requestWrapper struct {
	Body   []byte      `cbor:"body"`
	Method string      `cbor:"method"`
	URL    string      `cbor:"url"`
	Header http.Header `cbor:"header"`
}

type responseWrapper struct {
	Body       []byte      `cbor:"body"`
	StatusCode int         `cbor:"statusCode"`
	Header     http.Header `cbor:"header"`
}

type ResponseWritter struct {
	buffer     bytes.Buffer
	statusCode int
	header     http.Header
}

func (res *ResponseWritter) Header() http.Header {
	return res.header
}

func (res *ResponseWritter) Write(b []byte) (n int, err error) {
	return res.buffer.Write(b)
}

func (res *ResponseWritter) WriteHeader(status int) {
	res.statusCode = status
}
