package models_test

import (
	"booker/models"
	"database/sql"
	"testing"
)

func initTestingDB() *sql.DB {
	db := models.ConnectToDatabase("testing.db")
	models.CreateNewTables(db)
	models.FillWithSampleData(db)
	return db
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
}

func TestDate(t *testing.T) {
	db := initTestingDB()
	dates, err := models.GetDatesWithNamesNotBooked(db)
	checkError(t, err)

	date := dates[0]
	if date.BookedBy != -1 {
		t.Errorf("unbooked date.BookedBy = %d", date.BookedBy)
	}

	dateId := date.Id
	const userId = 123
	checkError(t, models.SetDateBookedBy(db, dateId, userId))

	newDate, err := models.GetDateById(db, dateId)
	checkError(t, err)
	if newDate.BookedBy != userId {
		t.Errorf("date booked by %d .BookedBy = %d", userId, newDate.BookedBy)
	}
	if newDate.Id != dateId {
		t.Errorf("received invalid date")
	}

	var bookedDates []*models.Date
	bookedDates, err = models.GetDatesBookedBy(db, userId)
	if len(bookedDates) != 1 {
		t.Errorf("received too many dates")
	}
	if bookedDates[0].Id != dateId {
		t.Errorf("received invalid date")
	}
}
