package main

import (
	"booker/http"
)

func main() {
	http.NewServer("booker.db").Run(":8080")
}
