package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	//This imports postgress driver, still don't understand how
	_ "github.com/lib/pq"
)

type ServerInfo struct {
	Address  string `json:"address"`
	SSLGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

type DomainInfo struct {
	IsDown      bool         `json:"is_down"`
	Severs      []ServerInfo `json:"servers"`
	SSLGrade    string       `json:"ssl_grade"`
	Title       string       `json:"title"`
	Logo        string       `json:"logo"`
	LastUpdated time.Time    `json:"last_updated"`
}

type DomainInfoRepo struct {
	Get    func(string) (DomainInfo, error)
	GetAll func() []DomainInfo
	Upsert func(DomainInfo)
}

func initDB(db *sql.DB) error {
	//Remember to
	db.Exec(`
	create table if not exists domains(
		host varchar not null primary key,
		data jsonb not null,
		last_updated timestamp not null
	);
	`)
	if _, err := db.Exec(`upsert into domains(host,data,last_updated)values($1,$2,$3)`, "pepe.com", "{}", time.Now()); err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query(`
		select * from domains;
	`)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		for rows.Next() {
			var host string
			var data string
			var lastUpdated time.Time
			err := rows.Scan(&host, &data, &lastUpdated)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("alskjdlakjsdjlkas")
				fmt.Println(host)
				fmt.Println(data)
				fmt.Println(lastUpdated)
			}
		}
		return nil
	}
}

func CreateRepo(dbURL string) (*DomainInfoRepo, error) {

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("aljksjkdlaskjdlaksjdlaskjd111")
		log.Fatal("error connecting to the database: ", err)
		return nil, err
	}
	initDB(db)

	return &DomainInfoRepo{
		Get: func(url string) (DomainInfo, error) {

			return DomainInfo{}, errors.New("No info for domain: " + url)
		},
		GetAll: func() []DomainInfo {

			return make([]DomainInfo, 0)
		},
		Upsert: func(seed DomainInfo) {

		},
	}, nil
}
