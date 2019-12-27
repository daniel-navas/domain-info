package storage

import (
	"database/sql"
	"encoding/json"
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
	Host        string       `json:"host"`
	IsDown      bool         `json:"is_down"`
	Severs      []ServerInfo `json:"servers"`
	SSLGrade    string       `json:"ssl_grade"`
	Title       string       `json:"title"`
	Logo        string       `json:"logo"`
	LastUpdated int          `json:"last_updated"`
}

// DomainInfoRepo :
type DomainInfoRepo struct {
	Get    func(string) (DomainInfo, error)
	GetAll func() ([]DomainInfo, error)
	Upsert func(DomainInfo) error
}

func initDB(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS domains (
			host VARCHAR NOT NULL PRIMARY KEY,
			data JSON NOT NULL,
			last_updated BIGINT NOT NULL
		);
	`); err != nil {
		log.Fatalln("Error:", err)
		return err
	}
	return nil
}

//CreateRepo :
func CreateRepo(dbURL string) (*DomainInfoRepo, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	initDB(db)

	return &DomainInfoRepo{
		Get: func(url string) (DomainInfo, error) {
			rows, err := db.Query(`SELECT d.host,d.data,d.last_updated FROM domains d WHERE d.host = $1	LIMIT 1`, url)
			if err != nil {
				return DomainInfo{}, err
			}
			defer rows.Close()
			for rows.Next() {
				var host string
				var data string
				var lastUpdate int
				rows.Scan(&host, &data, &lastUpdate)
				var domRecord DomainInfo
				json.Unmarshal([]byte(data), &domRecord)
				domRecord.LastUpdated = lastUpdate
				return domRecord, nil
			}
			return DomainInfo{}, errors.New("No records for domain: " + url)
		},
		GetAll: func() ([]DomainInfo, error) {
			records := make([]DomainInfo, 0)
			rows, err := db.Query(`SELECT d.host,d.data,d.last_updated FROM domains d`)
			if err != nil {
				return records, err
			}
			defer rows.Close()
			for rows.Next() {
				var host string
				var data string
				var lastUpdate int
				rows.Scan(&host, &data, &lastUpdate)
				var domRecord DomainInfo
				json.Unmarshal([]byte(data), &domRecord)
				domRecord.LastUpdated = lastUpdate
				records = append(records, domRecord)
			}
			if len(records) > 0 {
				return records, nil
			}
			return records, errors.New("No records")
		},
		Upsert: func(seed DomainInfo) error {
			seedJSON, err := json.Marshal(seed)
			if err != nil {
				log.Println("Error:", err)
				return err
			} else if _, err := db.Exec(
				`UPSERT INTO domains (host, data, last_updated) VALUES ($1, $2, $3)`,
				seed.Host, seedJSON, time.Now().Nanosecond()); err != nil {
				log.Fatalln("Error:", err)
				return err
			}
			return nil
		},
	}, nil
}
