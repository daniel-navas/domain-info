package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// AddressInfo :
type AddressInfo struct {
	CountryCode string `json:"countryCode"`
	ISP         string `json:"isp"`
}

// AddressInfoXtractor :
type AddressInfoXtractor struct{}

// Get :
func (aix *AddressInfoXtractor) Get(address string) AddressInfo {
	resp, err := http.Get("http://ip-api.com/json/" + address)
	if err != nil {
		log.Println("Error:", err)
		return AddressInfo{"unavailable", "unavailable"}
	}
	defer resp.Body.Close()
	var rawAddress AddressInfo
	json.NewDecoder(resp.Body).Decode(&rawAddress)

	return rawAddress
}
