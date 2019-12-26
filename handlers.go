package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anaskhan96/soup"
	"github.com/go-chi/chi"
)

func getDomainInfo(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + url)
	if err != nil {
		log.Printf("HTTP request failed. %s\n", err)
	}
	defer resp.Body.Close()
	var rawDom rawDomainInfo
	json.NewDecoder(resp.Body).Decode(&rawDom)
	domain := mapDomainInfo(rawDom)
	for index, server := range domain.Severs {
		resp, err := http.Get("http://ip-api.com/json/" + server.Address)
		if err != nil {
			log.Printf("HTTP request failed. %s\n", err)
		}
		defer resp.Body.Close()
		var rawAddress rawAddressInfo
		json.NewDecoder(resp.Body).Decode(&rawAddress)
		domain.Severs[index] = mapAddressInfo(rawAddress, server)
	}
	fmt.Println(getTitleAndLogo(url))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain)
	w.WriteHeader(http.StatusOK)
}

func getTitleAndLogo(url string) [2]string {
	resp, err := soup.Get("https://www.truora.com/")
	if err != nil {
		log.Printf("HTTP request failed. %s\n", err)
	}
	doc := soup.HTMLParse(resp)
	title := doc.Find("title")
	logo := doc.Find("link", "rel", "shortcut")
	var titleAndLogo [2]string
	if title.Error == nil {
		titleAndLogo[0] = title.Text()
	} else {
		log.Println(title.Error)
	}
	if logo.Error == nil {
		titleAndLogo[1] = logo.Attrs()["href"]
	} else {
		log.Println(title.Error)
	}
	return titleAndLogo
}
