package models

import (
	"database/sql"
	"time"
)

const sqlDateTable = `
DROP TABLE IF EXISTS dates;
CREATE TABLE dates (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	startTime  INTEGER NOT NULL,
	endTime	   INTEGER NOT NULL,
	bookedBy   INTEGER,
	assignedTo INTEGER NOT NULL,
	FOREIGN KEY(bookedBy) REFERENCES USER(id),
	FOREIGN KEY(assignedTo) REFERENCES USER(id)
);`

type Date struct {
	Id         int
	StartTime  time.Time
	EndTime    time.Time
	BookedBy   int // NULL is translated to -1
	AssignedTo int
}

func dateFromRow(row scannable) (*Date, error) {
	var u Date
	var start, end int64
	var bookedBy sql.NullInt32
	err := row.Scan(&u.Id, &start, &end, &bookedBy, &u.AssignedTo)
	u.StartTime = time.Unix(start, 0)
	u.EndTime = time.Unix(end, 0)
	u.BookedBy = int(bookedBy.Int32)
	if !bookedBy.Valid {
		u.BookedBy = -1
	}
	return &u, err
}

type DateWithNames struct {
	Date
	BookedByName   string
	AssignedToName string
}

func dateUserNamesFromRow(row scannable) (*DateWithNames, error) {
	var u DateWithNames
	var start, end int64
	var bookedBy sql.NullInt32
	err := row.Scan(&u.Id, &start, &end, &bookedBy, &u.AssignedTo, &u.BookedByName, &u.AssignedToName)
	u.StartTime = time.Unix(start, 0)
	u.EndTime = time.Unix(end, 0)
	u.BookedBy = int(bookedBy.Int32)
	if !bookedBy.Valid {
		u.BookedBy = -1
	}
	return &u, err
}

const sqlDateAllNotBooked = `
SELECT dates.*, '', emp.name
FROM dates 
LEFT JOIN users emp ON dates.assignedTo = emp.id
WHERE bookedBy IS NULL`

func GetDatesWithNamesNotBooked(db *sql.DB) ([]*DateWithNames, error) {
	rows, err := db.Query(sqlDateAllNotBooked)
	if err != nil {
		return nil, err
	}
	return readFromRows(rows, dateUserNamesFromRow)
}

const sqlDateCreate = `
INSERT INTO dates (startTime, endTime, assignedTo) VALUES (?, ?, ?)`

func CreateDate(db *sql.DB, startTime time.Time, endTime time.Time, assignedTo int) error {
	_, err := db.Exec(sqlDateCreate, startTime.Unix(), endTime.Unix(), assignedTo)
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

const sqlDateAssignedTo = `
SELECT dates.*, IFNULL(cus.name, ''), emp.name
FROM dates 
LEFT JOIN users cus ON dates.bookedBy = cus.Id
LEFT JOIN users emp ON dates.assignedTo = emp.Id
WHERE assignedTo = ?`

func GetDatesWithNamesAssignedTo(db *sql.DB, empId int) ([]*DateWithNames, error) {
	rows, err := db.Query(sqlDateAssignedTo, empId)
	if err != nil {
		return nil, err
	}
	return readFromRows(rows, dateUserNamesFromRow)
}

const sqlDateAll = `
SELECT dates.*, IFNULL(cus.name, ''), emp.name
FROM dates 
LEFT JOIN users cus ON dates.bookedBy = cus.Id
LEFT JOIN users emp ON dates.assignedTo = emp.Id`

func GetDatesWithNamesAll(db *sql.DB) ([]*DateWithNames, error) {
	rows, err := db.Query(sqlDateAll)
	if err != nil {
		return nil, err
	}
	return readFromRows(rows, dateUserNamesFromRow)
}
