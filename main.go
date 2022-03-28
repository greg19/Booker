package main

import (
	"booker/http"
)

func main() {
	http.NewServer().Run(":8080")
}
