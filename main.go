package main

import (
	"log"
	"net/http"

	"github.com/dfnavas/domain-info/controllers"
	"github.com/dfnavas/domain-info/storage"

	"github.com/dfnavas/domain-info/middleware"
)

func main() {

	tAiXtractor := middleware.CreateTitleAndLogoXtractor()

	domainInfoXtrator := middleware.CreateDomainInfoXtractor()

	repo, err := storage.CreateRepo("postgresql://root@localhost:26257?sslmode=disable")

	if err != nil {
		log.Fatal("Error creating repo")
	} else {

		addressInfoXtractor := middleware.CreateAddressInfoXtractor()

		ctrl := controllers.CreateCtrl(tAiXtractor, domainInfoXtrator, addressInfoXtractor, repo)
		router := createRouter(ctrl)
		http.ListenAndServe(":3000", router)
	}
}
