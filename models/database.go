package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectToDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "booker.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateNewTables(db *sql.DB) {
	tables := []string{
		sqlDateTable,
		sqlUserTable,
		sqlSessionTable,
	}

	for _, query := range tables {
		if _, err := db.Exec(query); err != nil {
			log.Print("query:", query, "\n")
			log.Fatal(err)
		}
	}
}

type scannable interface {
	Scan(...interface{}) error
}

func readFromRows[T any](rows *sql.Rows, callback func(scannable) (*T, error)) ([]*T, error) {
	defer rows.Close()
	var items []*T
	for rows.Next() {
		i, err := callback(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

