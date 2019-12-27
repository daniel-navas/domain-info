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

// DomainHistory :
type DomainHistory struct {
	Host     string       `json:"host"`
	Severs   []ServerInfo `json:"servers"`
	SSLGrade string       `json:"ssl_grade"`
	Title    string       `json:"title"`
	Logo     string       `json:"logo"`
}

// DomainInfoRepo :
type DomainInfoRepo struct {
	Get    func(string) (DomainInfo, error)
	GetAll func() []DomainHistory
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

// CreateRepo :
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
		GetAll: func() []DomainHistory {
			records := make([]DomainHistory, 0)
			rows, err := db.Query(`SELECT d.host,d.data FROM domains d`)
			if err != nil {
				log.Println("Error:", err)
				return records
			}
			defer rows.Close()
			for rows.Next() {
				var host string
				var data string
				rows.Scan(&host, &data)
				var domRecord DomainHistory
				json.Unmarshal([]byte(data), &domRecord)
				records = append(records, domRecord)
			}
			return records
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
