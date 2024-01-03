package sdk

import (
	"bytes"
	"net/http"
	"unsafe"

	"github.com/fxamacker/cbor/v2"
)

type Point struct {
	X int
	Y int
}

//go:wasmimport env getEncodedRequestSize
//go:noescape
func getEncodedRequestSize() uint32

//go:wasmimport env readEncodedRequestToPointer
//go:noescape
func readEncodedRequestToPointer(pointer uint32)

//go:wasmimport env writeResponseFromPointer
//go:noescape
func writeResponseFromPointer(pointer uint32, size uint32)

func HandleFunc(callback http.HandlerFunc) {
	encodedRequestSize := getEncodedRequestSize()
	encodedRequestBuffer := make([]byte, encodedRequestSize)

	readEncodedRequestToPointer(uint32(uintptr(unsafe.Pointer(&encodedRequestBuffer[0]))))

	var reqWrapper requestWrapper
	if err := cbor.Unmarshal(encodedRequestBuffer, &reqWrapper); err != nil {
	}

	httpReq, _ := http.NewRequest(reqWrapper.Method, reqWrapper.URL, bytes.NewReader(reqWrapper.Body))
	httpReq.Header = reqWrapper.Header

	res := &ResponseWritter{
		header: http.Header{},
	}
	callback(res, httpReq)

	resEncodedBuf, _ := cbor.Marshal(responseWrapper{
		Body:       res.buffer.Bytes(),
		StatusCode: res.statusCode,
		Header:     res.header,
	})
	writeResponseFromPointer(uint32(uintptr(unsafe.Pointer(&resEncodedBuf[0]))), uint32(len(resEncodedBuf)))
}
