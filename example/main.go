package main

import (
	"fmt"
	"net/http"

	"github.com/Stolkerve/kappa/sdk"
)

func main() {
	sdk.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("User-Agent"))
		w.WriteHeader(201)
		w.Header().Add("pepe", "123")
		w.Write([]byte("hola mundo :O"))
	})
}
