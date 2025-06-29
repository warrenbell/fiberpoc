package main

import (
	"flag"
	"log"
	"os"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error: B2WRUD - Loading .env. Error: %v", err)
	}

	migrationDirection := flag.String("direct", "up", "The migration direction, up or down")
	flag.Parse()

	log.Printf("Info: BP04Y1 - Migration direction is %s", *migrationDirection)

	m, err := migrate.New(
		"file://./migrations",
		os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Printf("Error: OUCUGO - Creating new migration. Error: %v", err)
	}
	if *migrationDirection == "up" {
		if err := m.Up(); err != nil {
			log.Printf("Error: 99ECW0 - Running up migration. Error: %v", err)
		}
	} else {
		if err := m.Down(); err != nil {
			log.Printf("Error: OF1VW6 - Running down migration. Error: %v", err)
		}
	}
}
