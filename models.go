package main

type rawDomainInfo struct {
	Status string `json:"status"`
}

type domainInfo struct {
	IsDown bool `json:"is_down"`
}
