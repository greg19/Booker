package main

import (
	"booker/models"
)

func main() {
	db := models.ConnectToDatabase("booker.db")
	models.CreateNewTables(db)
	models.FillWithSampleData(db)
}
