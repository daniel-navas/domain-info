package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// DomainInfoXtractor :
type DomainInfoXtractor interface {
	Get(url string) (*DomainInfo, error)
}

type domainInfoXtractorDeps struct{}

// DomainInfo :
type DomainInfo struct {
	Host      string         `json:"host"`
	Status    string         `json:"status"`
	Endpoints []endpointInfo `json:"endpoints"`
}

type endpointInfo struct {
	IPAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
}

// NewDomainInfoXtractor :
func NewDomainInfoXtractor() DomainInfoXtractor {
	return &domainInfoXtractorDeps{}
}

// Get :
func (dix *domainInfoXtractorDeps) Get(url string) (*DomainInfo, error) {
	resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + url)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	var rawDom DomainInfo
	json.NewDecoder(resp.Body).Decode(&rawDom)
	return &rawDom, nil
}
