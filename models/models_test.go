package models_test

import (
	"booker/models"
	"database/sql"
	"testing"
	"time"
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

func checkUserType(t *testing.T, u *models.User, expectedUserType int) {
	if u.IsAdmin() != (expectedUserType == models.UserTypeAdmin) {
		t.Errorf("u.IsAdmin: unexpected result")
	}
	if u.IsEmployee() != (expectedUserType <= models.UserTypeEmployee) {
		t.Errorf("u.IsEmployee: unexpected result")
	}
	if u.IsCustomer() != true {
		t.Errorf("u.IsCustomer: unexpected result")
	}
}

func checkArraySize[T any](t *testing.T, arr []T, expectedSize int) {
	if len(arr) != expectedSize {
		t.Errorf("array length: %d expected: %d", len(arr), expectedSize)
	}
}

func TestUser(t *testing.T) {
	const adminId = 1
	const employeeId = 2
	const customerId = 4

	db := initTestingDB()

	admin, err := models.GetUserById(db, adminId)
	checkError(t, err)
	checkUserType(t, admin, models.UserTypeAdmin)

	employee, err := models.GetUserById(db, employeeId)
	checkError(t, err)
	checkUserType(t, employee, models.UserTypeEmployee)

	customer, err := models.GetUserById(db, customerId)
	checkError(t, err)
	checkUserType(t, customer, models.UserTypeCustomer)

	admins, err := models.GetUsersByType(db, models.UserTypeAdmin)
	checkError(t, err)
	checkArraySize(t, admins, 1)
}

func TestSession(t *testing.T) {
	db := initTestingDB();

	const token = "secret"

	err := models.CreateSession(db, token, 1, time.Now().Add(time.Hour))
	checkError(t, err)

	sess, err := models.GetSessionByToken(db, token)
	checkError(t, err)

	if sess.IsExpired() {
		t.Errorf("session should not be expired");
	}

	err = models.DeleteSession(db, token)
	checkError(t, err)

/*
	sess, err = models.GetSessionByToken(db, token)
	checkError(t, err)

	if sess != nil {
		t.Errorf("session should be deleted")
	}
*/
}
