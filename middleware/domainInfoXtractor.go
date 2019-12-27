package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type EndpointInfo struct {
	IPAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
}
type DomainInfo struct {
	Status    string         `json:"status"`
	Endpoints []EndpointInfo `json:"endpoints"`
}

type DomainInfoXtractor struct {
	Get func(string) DomainInfo
}

func CreateDomainInfoXtractor() *DomainInfoXtractor {
	return &DomainInfoXtractor{
		Get: func(url string) DomainInfo {
			//TODO

			resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + url)
			if err != nil {
				log.Printf("HTTP request failed. %s\n", err)
			}
			defer resp.Body.Close()
			var rawDom DomainInfo
			json.NewDecoder(resp.Body).Decode(&rawDom)
			return rawDom
		},
	}
}
