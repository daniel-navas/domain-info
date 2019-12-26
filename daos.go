package main

func mapDomainInfo(rawDom rawDomainInfo) domainInfo {
	var grades = map[string]int{"A+": 10, "A": 9, "A-": 8, "B": 7, "C": 6, "D": 5, "E": 4, "F": 3, "M": 2, "T": 1}
	var domain domainInfo
	lowestGrade := "A+"
	domain.IsDown = rawDom.Status == "DNS"
	for _, ep := range rawDom.Endpoints {
		var server serverInfo
		server.Address = ep.IPAddress
		server.SSLGrade = ep.Grade
		domain.Severs = append(domain.Severs, server)
		if grades[ep.Grade] < grades[lowestGrade] {
			lowestGrade = ep.Grade
		}
	}
	domain.SSLGrade = lowestGrade
	return domain
}

func mapAddressInfo(rawAddress rawAddressInfo, server serverInfo) serverInfo {
	server.Country = rawAddress.CountryCode
	server.Owner = rawAddress.ISP
	return server
}

type rawDomainInfo struct {
	Status    string         `json:"status"`
	Endpoints []endpointInfo `json:"endpoints"`
}

type endpointInfo struct {
	IPAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
}

type rawAddressInfo struct {
	CountryCode string `json:"countryCode"`
	ISP         string `json:"isp"`
}

type domainInfo struct {
	IsDown   bool         `json:"is_down"`
	Severs   []serverInfo `json:"servers"`
	SSLGrade string       `json:"ssl_grade"`
}

type serverInfo struct {
	Address  string `json:"address"`
	SSLGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}
