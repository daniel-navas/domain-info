package controllers

import (
	"fmt"

	"github.com/dfnavas/domain-info/middleware"
	"github.com/dfnavas/domain-info/storage"
)

type SearchHistory struct {
	items []DomainInfo
}

type DomainInfo struct {
	IsDown           bool                 `json:"is_down"`
	Severs           []storage.ServerInfo `json:"servers"`
	SSLGrade         string               `json:"ssl_grade"`
	PreviousSSLGrade string               `json:"previous_ssl_grade"`
	ServersChanged   bool                 `json:"servers_changed"`
	Title            string               `json:"title"`
	Logo             string               `json:"logo"`
}

type DomainInfoCtrl struct {
	Get func(string) DomainInfo

	GetAll func() SearchHistory
}

func mapDomainInfo(rawDom middleware.DomainInfo,
	addressesInfo []middleware.AddressInfo,
	titleAndLogo middleware.TitleAndLogo) storage.DomainInfo {
	var grades = map[string]int{"A+": 10, "A": 9, "A-": 8, "B": 7, "C": 6, "D": 5, "E": 4, "F": 3, "M": 2, "T": 1}
	var domain storage.DomainInfo
	lowestGrade := "A+"
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
	domain.Title = titleAndLogo.Title
	domain.Logo = titleAndLogo.Logo
	return domain
}

func CreateCtrl(
	tAiXtractor *middleware.TitleAndLogoXtractor,
	infoXtractor *middleware.DomainInfoXtractor,
	addressInfoXtractor *middleware.AddressInfoXtractor,
	repo *storage.DomainInfoRepo) *DomainInfoCtrl {
	return &DomainInfoCtrl{
		Get: func(url string) DomainInfo {
			info := infoXtractor.Get(url)
			totalEndpoints := len(info.Endpoints)
			servers := make([]middleware.AddressInfo, totalEndpoints)
			for idx, endpoint := range info.Endpoints {
				servers[idx] = addressInfoXtractor.Get(endpoint.IPAddress)
			}
			titleAndLogo := tAiXtractor.Get(url)
			var currentInfo storage.DomainInfo = mapDomainInfo(info, servers, titleAndLogo)
			previousInfo, err := repo.Get(url)
			repo.Upsert(currentInfo)
			result := DomainInfo{
				IsDown:   currentInfo.IsDown,
				Severs:   currentInfo.Severs,
				SSLGrade: currentInfo.SSLGrade,
				Title:    currentInfo.Title,
				Logo:     currentInfo.Logo,
			}
			if err != nil {
				return result
			} else {
				fmt.Println(previousInfo)
				//TODO add PreviousSSLGrade and ServersChanged
				return result
			}
		},
		GetAll: func() SearchHistory {
			return SearchHistory{
				items: make([]DomainInfo, 0),
			}
		},
	}
}
