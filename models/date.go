package models

import (
	"database/sql"
	"time"
)

const sqlDateTable = `
DROP TABLE IF EXISTS dates;
CREATE TABLE dates (
	id        INTEGER PRIMARY KEY AUTOINCREMENT,
	startTime INTEGER NOT NULL,
	endTime	  INTEGER NOT NULL,
	bookedBy  INTEGER
);`

type Date struct {
	Id        int
	StartTime time.Time
	EndTime   time.Time
	BookedBy  int // NULL is translated to -1
}

const sqlDateAllNotBooked = `
SELECT * FROM dates WHERE bookedBy IS NULL`

func dateFromRow(row scannable) (*Date, error) {
	var u Date
	var start, end int64
	var bookedBy sql.NullInt32
	err := row.Scan(&u.Id, &start, &end, &bookedBy)
	u.StartTime = time.Unix(start, 0)
	u.EndTime = time.Unix(end, 0)
	u.BookedBy = int(bookedBy.Int32)
	if !bookedBy.Valid {
		u.BookedBy = -1
	}
	return &u, err
}

func GetDatesNotBooked(db *sql.DB) ([]*Date, error) {
	rows, err := db.Query(sqlDateAllNotBooked)
	if err != nil {
		return nil, err
	}
	return readFromRows(rows, dateFromRow)
}

const sqlDateCreate = `
INSERT INTO dates (startTime, endTime) VALUES (?, ?)`

func CreateDate(db *sql.DB, startTime time.Time, endTime time.Time) error {
	_, err := db.Exec(sqlDateCreate, startTime.Unix(), endTime.Unix())
	return err
}

const sqlDateBookedBy = `
SELECT * FROM dates WHERE bookedBy = ?`

func GetDatesBookedBy(db *sql.DB, userId int) ([]*Date, error) {
	rows, err := db.Query(sqlDateBookedBy, userId)
	if err != nil {
		return nil, err
	}
	return readFromRows(rows, dateFromRow)
}

const sqlDateById = `
SELECT * FROM dates WHERE id = ?`

func GetDateById(db *sql.DB, id int) (*Date, error) {
	row := db.QueryRow(sqlDateById, id)
	return dateFromRow(row)
}

const sqlDateSetBookedBy = `
UPDATE dates SET bookedBy = ? WHERE id = ?`

func SetDateBookedBy(db *sql.DB, dateId int, userId int) error {
	var err error
	if userId != -1 {
		_, err = db.Exec(sqlDateSetBookedBy, userId, dateId)
	} else {
		_, err = db.Exec(sqlDateSetBookedBy, nil, dateId)
	}
	return err
}
