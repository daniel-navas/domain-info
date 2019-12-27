package storage

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

// ServerInfo :
type ServerInfo struct {
	Address  string `json:"address"`
	SSLGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

// DomainInfo :
type DomainInfo struct {
	IsDown      bool         `json:"is_down"`
	Severs      []ServerInfo `json:"servers"`
	SSLGrade    string       `json:"ssl_grade"`
	Title       string       `json:"title"`
	Logo        string       `json:"logo"`
	LastUpdated time.Time    `json:"last_updated"`
}

// DomainInfoRepo :
type DomainInfoRepo struct {
	Get    func(string) (DomainInfo, error)
	GetAll func() []DomainInfo
	Upsert func(DomainInfo)
}

func initDB(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS domains (
			host VARCHAR NOT NULL PRIMARY KEY,
			data JSON NOT NULL,
			last_update TIMESTAMP NOT NULL
		);
	`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		`UPSERT INTO domains (host, data, last_update) VALUES ($1, $2, $3)`, "test.com", "{}", time.Now()); err != nil {
		log.Fatal(err)
	}

	// rows, err := db.Query(`
	// 	select * from domains;
	// `)

	// defer rows.Close()

	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// for rows.Next() {
	// 	var host string
	// 	var data string
	// 	var lastUpdated time.Time
	// 	err := rows.Scan(&host, &data, &lastUpdated)

	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Println(host)
	// 		fmt.Println(data)
	// 		fmt.Println(lastUpdated)
	// 	}
	// }

	return nil
}

//CreateRepo :
func CreateRepo(dbURL string) (*DomainInfoRepo, error) {
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
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
