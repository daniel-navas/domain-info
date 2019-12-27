package controllers

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dfnavas/domain-info/middleware"
	"github.com/dfnavas/domain-info/storage"
)

type searchHistory struct {
	items []DomainInfo
}

// DomainInfo :
type DomainInfo struct {
	IsDown           bool                 `json:"is_down"`
	Severs           []storage.ServerInfo `json:"servers"`
	SSLGrade         string               `json:"ssl_grade"`
	PreviousSSLGrade string               `json:"previous_ssl_grade"`
	ServersChanged   bool                 `json:"servers_changed"`
	Title            string               `json:"title"`
	Logo             string               `json:"logo"`
}

// DomainInfoCtrl :
type DomainInfoCtrl struct {
	Get    func(string) (DomainInfo, error)
	GetAll func() searchHistory
}

func mapDomainInfo(rawDom middleware.DomainInfo,
	addressesInfo []middleware.AddressInfo,
	title string, logo string) storage.DomainInfo {
	var domain storage.DomainInfo
	var grades = map[string]int{"A+": 10, "A": 9, "A-": 8, "B": 7, "C": 6, "D": 5, "E": 4, "F": 3, "M": 2, "T": 1}
	lowestGrade := "A+"
	domain.Host = rawDom.Host
	domain.IsDown = rawDom.Status == "DNS"
	for idx, ep := range rawDom.Endpoints {
		var server storage.ServerInfo
		server.Address = ep.IPAddress
		server.SSLGrade = ep.Grade
		server.Owner = addressesInfo[idx].ISP
		server.Country = addressesInfo[idx].CountryCode
		domain.Severs = append(domain.Severs, server)
		if grades[ep.Grade] < grades[lowestGrade] {
			lowestGrade = ep.Grade
		}
	}
	domain.SSLGrade = lowestGrade
	domain.Title = title
	domain.Logo = logo
	return domain
}

// CreateCtrl :
func CreateCtrl(
	tagXtractor *middleware.TagXtractor,
	domainInfoXtractor *middleware.DomainInfoXtractor,
	addressInfoXtractor *middleware.AddressInfoXtractor,
	repo *storage.DomainInfoRepo) *DomainInfoCtrl {
	return &DomainInfoCtrl{
		Get: func(url string) (DomainInfo, error) {
			//FIRST WE JOIN ALL INFO-------
			// Domain Info
			rawDomInfo, err := domainInfoXtractor.Get(url)
			if err != nil {
				return DomainInfo{}, err
			}
			//Then we get info per endpoint/server of domain
			endpointsLength := len(rawDomInfo.Endpoints)
			servers := make([]middleware.AddressInfo, endpointsLength) //TODO change slice for array?
			for idx, endpoint := range rawDomInfo.Endpoints {
				servers[idx] = addressInfoXtractor.Get(endpoint.IPAddress)
			}
			//Then we get title and logo info
			title, logo := tagXtractor.Get(url)
			//And then we aggregate it all
			var newDomInfo storage.DomainInfo = mapDomainInfo(rawDomInfo, servers, title, logo)
			//Then we get a previous registry of the given URL
			oldDomInfo, err := repo.Get(url)
			//Then we save the new one
			repo.Upsert(newDomInfo)
			//Finally, if there was no previous registry, or it was TOO old, then we ignore it and return
			// Without providing PreviousSSLGrade or ServersChanged
			fmt.Println("qqqq", oldDomInfo.LastUpdated)
			fmt.Println("wwww", time.Now().Add(time.Duration(-3.6e+12)).Nanosecond())
			fmt.Println("eeee", oldDomInfo.LastUpdated < time.Now().Add(time.Duration(-3.6e+12)).Nanosecond())
			if err != nil || oldDomInfo.LastUpdated < time.Now().Add(time.Duration(-3.6e+12)).Nanosecond() {
				return DomainInfo{
					IsDown:   newDomInfo.IsDown,
					Severs:   newDomInfo.Severs,
					SSLGrade: newDomInfo.SSLGrade,
					Title:    newDomInfo.Title,
					Logo:     newDomInfo.Logo,
				}, nil
			}
			//Otherwise we return the full info
			return DomainInfo{
				IsDown:           newDomInfo.IsDown,
				Severs:           newDomInfo.Severs,
				SSLGrade:         newDomInfo.SSLGrade,
				Title:            newDomInfo.Title,
				Logo:             newDomInfo.Logo,
				PreviousSSLGrade: oldDomInfo.SSLGrade,
				ServersChanged:   !reflect.DeepEqual(newDomInfo.Severs, oldDomInfo.Severs),
			}, nil

		},
		GetAll: func() searchHistory {
			return searchHistory{
				items: make([]DomainInfo, 0),
			}
		},
	}
}
