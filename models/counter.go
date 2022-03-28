package models

// example model

import (
	"database/sql"
	"log"
)

type Counter struct {
	Name string
	Cnt  int
}

func InitTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS counters (name TEXT PRIMARY KEY, cnt INT)")
	if err != nil {
		log.Fatal(err)
	}
}

func GetCounter(db *sql.DB, name string) *Counter {
	rows, err := db.Query("SELECT * FROM counters where name = ?", name)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var counter Counter
		err := rows.Scan(&counter.Name, &counter.Cnt)
		if err != nil {
			log.Println(err)
		}
		return &counter
	}

	return nil
}

func UpdateCounter(db *sql.DB, name string, cnt int) {
	_, err := db.Exec("UPDATE counters SET cnt = ? WHERE name = ?", cnt, name)
	if err != nil {
		log.Fatal(err)
	}
}

func AddCounter(db *sql.DB, counter *Counter) {
	_, err := db.Exec("INSERT INTO counters (name, cnt) VALUES (?, ?)", counter.Name, counter.Cnt)
	if err != nil {
		log.Fatal(err)
	}
}
