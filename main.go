package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{url}", getDomainInfo)
	http.ListenAndServe(":3000", r)
}
