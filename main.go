package main

import (
	_ "github.com/lib/pq"

	"challenge-goapi/config"
)

var db = config.ConnectDB()

func main() {
	// Tulis kode kamu disini
	// db := connectDB()
	defer db.Close()
}
