package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type endpointInfo struct {
	IPAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
}

// DomainInfo :
type DomainInfo struct {
	Host      string         `json:"host"`
	Status    string         `json:"status"`
	Endpoints []endpointInfo `json:"endpoints"`
}

// DomainInfoXtractor :
type DomainInfoXtractor struct {
	Get func(string) (DomainInfo, error)
}

// CreateDomainInfoXtractor :
func CreateDomainInfoXtractor() *DomainInfoXtractor {
	return &DomainInfoXtractor{
		Get: func(url string) (DomainInfo, error) {
			var rawDom DomainInfo
			resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + url) //TODO what should I return if this fails
			if err != nil {
				log.Printf("HTTP request failed. %s\n", err)
				return rawDom, err
			}
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&rawDom)

			return rawDom, nil
		},
	}
}
