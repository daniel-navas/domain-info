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
type AddressInfoXtractor struct {
	Get func(string) AddressInfo
}

// CreateAddressInfoXtractor :
func CreateAddressInfoXtractor() *AddressInfoXtractor {

	return &AddressInfoXtractor{
		Get: func(address string) AddressInfo {
			resp, err := http.Get("http://ip-api.com/json/" + address)
			if err != nil {
				log.Printf("HTTP request failed. %s\n", err)
			}
			defer resp.Body.Close()
			var rawAddress AddressInfo
			json.NewDecoder(resp.Body).Decode(&rawAddress)

			return rawAddress
		},
	}

}
