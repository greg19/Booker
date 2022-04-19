package main

import (
	"booker/models"
	"database/sql"
	"log"
	"time"
)

func fillWithSampleData(db *sql.DB) {
	users := []models.User{
		{Name: "Admin", Username: "admin", Password: "admin", UserType: models.UserTypeAdmin},
		{Name: "Andrzej", Username: "pracownik", Password: "roku", UserType: models.UserTypeEmployee},
		{Name: "bob", Username: "bob", Password: "123", UserType: models.UserTypeCustomer},
	}
	for _, c := range users {
		if err := models.CreateUser(db, c.Name, c.Username, c.Password, c.UserType); err != nil {
			log.Fatal(err)
		}
	}

	for i := 1; i <= 5; i++ {
		startTime := time.Now().Add(time.Duration(i) * time.Hour)
		endTime := time.Now().Add(time.Duration(i+1) * time.Hour)
		if err := models.CreateDate(db, startTime, endTime); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	db := models.ConnectToDatabase()
	models.CreateNewTables(db)
	fillWithSampleData(db)
}
