package controllers

import (
	"reflect"
	"time"

	"github.com/dfnavas/domain-info/middleware"
	"github.com/dfnavas/domain-info/storage"
)

// DomainInfoCtrl :
type DomainInfoCtrl interface {
	Get(string) (*DomainInfo, error)
	GetAll() []*storage.DomainHistory
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

type domainInfoCtrlDeps struct {
	tagXtractor         middleware.TagXtractor
	domainInfoXtractor  middleware.DomainInfoXtractor
	addressInfoXtractor middleware.AddressInfoXtractor
	repo                storage.DomainInfoRepo
}

// CreateCtrl :
func CreateCtrl(
	tagXtractor middleware.TagXtractor,
	domainInfoXtractor middleware.DomainInfoXtractor,
	addressInfoXtractor middleware.AddressInfoXtractor,
	repo storage.DomainInfoRepo) DomainInfoCtrl {
	return &domainInfoCtrlDeps{tagXtractor,
		domainInfoXtractor,
		addressInfoXtractor,
		repo}
}

func (ctrl *domainInfoCtrlDeps) Get(url string) (*DomainInfo, error) {
	// Get domain info from api.ssllabs.com
	rawDomInfo, err := ctrl.domainInfoXtractor.Get(url)
	if err != nil {
		return nil, err
	}
	// Get address infor per server of domain
	endpointsLength := len(rawDomInfo.Endpoints)
	servers := make([]*middleware.AddressInfo, endpointsLength)
	for idx, endpoint := range rawDomInfo.Endpoints {
		info, err := ctrl.addressInfoXtractor.Get(endpoint.IPAddress)
		if err != nil {
			return nil, err
		}
		servers[idx] = info
	}
	// Get title and logo info
	title, logo := ctrl.tagXtractor.GetTitleAndLogo(url)
	// Map everything in one object
	var newDomInfo storage.DomainInfo = *mapDomainInfo(rawDomInfo, servers, title, logo)
	// Get a previous record of the given host
	oldDomInfo, err := ctrl.repo.Get(url)
	// Save the new one record
	ctrl.repo.Upsert(newDomInfo)
	// If no record or the last record is to young (1Hr)
	// Return the domain info without providing PreviousSSLGrade or ServersChanged
	domInfo := DomainInfo{
		IsDown:   newDomInfo.IsDown,
		Severs:   newDomInfo.Severs,
		SSLGrade: newDomInfo.SSLGrade,
		Title:    newDomInfo.Title,
		Logo:     newDomInfo.Logo,
	}
	if err != nil || oldDomInfo == nil || oldDomInfo.LastUpdated < time.Now().Add(time.Duration(-3.6e+12)).UnixNano() {
		return &domInfo, nil
	}
	// Otherwise return domain info with PreviousSSLGrade and ServersChanged
	domInfo.PreviousSSLGrade = oldDomInfo.SSLGrade
	domInfo.ServersChanged = !reflect.DeepEqual(newDomInfo.Severs, oldDomInfo.Severs)
	return &domInfo, nil
}

func (ctrl *domainInfoCtrlDeps) GetAll() []*storage.DomainHistory {
	history := ctrl.repo.GetAll()
	return history
}

func mapDomainInfo(rawDom *middleware.DomainInfo,
	addressesInfo []*middleware.AddressInfo,
	title string, logo string) *storage.DomainInfo {
	var domain storage.DomainInfo
	domain.Host = rawDom.Host
	domain.IsDown = rawDom.Status == "DNS" || rawDom.Status == "ERROR" || rawDom.Status == ""
	domain.Severs = []storage.ServerInfo{}
	if domain.IsDown {
		return &domain
	}
	var grades = map[string]int{"A+": 10, "A": 9, "A-": 8, "B": 7, "C": 6, "D": 5, "E": 4, "F": 3, "M": 2, "T": 1}
	lowestGrade := "A+"
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
	return &domain
}
