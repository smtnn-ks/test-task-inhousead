package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/smtnn-ks/test-task-inhousead/db"
	"github.com/smtnn-ks/test-task-inhousead/router"
	"github.com/smtnn-ks/test-task-inhousead/scraper"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	db.Init()
	defer db.Close()

	go scraper.Init()

	app := router.Init()

	log.Fatal(app.Listen(":" + port))
}
