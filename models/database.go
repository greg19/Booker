package models

import (
	"database/sql"
	"log"
	"time"
	"path"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectToDatabase(dbfilename string) *sql.DB {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	db, err := sql.Open("sqlite3", path.Join(basepath, "..", dbfilename))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type dbtype interface {
	Exec(string, ...interface{}) (sql.Result, error)	
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
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

func FillWithSampleData(db *sql.DB) {
	users := []User{
		{Name: "Admin", Username: "admin", Password: "admin", UserType: UserTypeAdmin},
		{Name: "Andrzej", Username: "pracownik", Password: "roku", UserType: UserTypeEmployee},
		{Name: "Fabian", Username: "pracownik2", Password: "miesiaca", UserType: UserTypeEmployee},
		{Name: "bob", Username: "bob", Password: "123", UserType: UserTypeCustomer},
	}
	for _, c := range users {
		if err := CreateUser(db, c.Name, c.Username, c.Password, c.UserType); err != nil {
			log.Fatal(err)
		}
	}

	for _, username := range []string {"pracownik", "pracownik2"} {
		emp, err := GetUserByUsername(db, username)
		if err != nil {
			log.Fatal(err)
		}

		for i := 1; i <= 5; i++ {
			startTime := time.Now().Add(time.Duration(i) * time.Hour)
			endTime := time.Now().Add(time.Duration(i+1) * time.Hour)
			if err := CreateDate(db, startTime, endTime, emp.Id); err != nil {
				log.Fatal(err)
			}
		}
	}
}
