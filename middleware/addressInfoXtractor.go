package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// AddressInfoXtractor :
type AddressInfoXtractor interface {
	Get(address string) (*AddressInfo, error)
}

type addressInfoXtractorDeps struct{}

// AddressInfo :
type AddressInfo struct {
	CountryCode string `json:"countryCode"`
	ISP         string `json:"isp"`
}

// NewAddressInfoXtractor :
func NewAddressInfoXtractor() AddressInfoXtractor {
	return &addressInfoXtractorDeps{}
}

// Get :
func (aix *addressInfoXtractorDeps) Get(address string) (*AddressInfo, error) {
	resp, err := http.Get("http://ip-api.com/json/" + address)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	var rawAddress AddressInfo
	json.NewDecoder(resp.Body).Decode(&rawAddress)
	return &rawAddress, nil
}
