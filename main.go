package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{url}", getDomainInfo)
	http.ListenAndServe(":3000", r)
}

func getDomainInfo(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + url)
	if err != nil {
		log.Fatalf("HTTP request failed. %s\n", err)
	}
	defer resp.Body.Close()
	var rawDom rawDomainInfo
	json.NewDecoder(resp.Body).Decode(&rawDom)
	domain := mapDomainInfo(rawDom)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
	w.WriteHeader(http.StatusOK)

}
