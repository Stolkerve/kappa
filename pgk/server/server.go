package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Stolkerve/kappa/pgk/api"
	"github.com/go-chi/chi/v5"
)

func NewServer(addrs string) {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		api.DeployRoutes(r)
		api.CallRoutes(r)
	})

	fmt.Println("Running in ", addrs)
	err := http.ListenAndServe(addrs, r)
	if err != nil {
		log.Fatalln(err)
	}
}
