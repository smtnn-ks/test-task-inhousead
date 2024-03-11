package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var Client *sql.DB

func Init() {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not specified")
	}

	var err error
	Client, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}
}

func Close() error {
	return Client.Close()
}
