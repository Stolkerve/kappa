#!/bin/bash
# GOOS=wasip1 GOARCH=wasm go build -o ../cmd/function.wasm main.go
tinygo build -o ../cmd/function.wasm -target=wasi -no-debug main.go