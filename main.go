package main

import (
	"log"
	"net/http"

	"github.com/dfnavas/domain-info/controllers"
	"github.com/dfnavas/domain-info/middleware"
	"github.com/dfnavas/domain-info/storage"
	_ "github.com/lib/pq"
)

func main() {
	domainInfoXtrator := middleware.CreateDomainInfoXtractor()
	tagXtractor := middleware.CreateTagXtractor()

	repo, err := storage.CreateRepo("postgresql://maxroach@localhost:26257/domainsdb?sslmode=disable")
	if err != nil {
		log.Fatalln("Error connecting to the database: ", err)
	} else {
		addressInfoXtractor := middleware.CreateAddressInfoXtractor()
		ctrl := controllers.CreateCtrl(tagXtractor, domainInfoXtrator, addressInfoXtractor, repo)
		router := createRouter(ctrl)
		http.ListenAndServe(":3333", router)
	}
}
