package main

import (
	"fmt"
	"log"

	"github.com/contextart/al/api/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	const migrationsFolder = "file://db/sql/migrations"

	log.Println("Running migrations script")

	settings := config.Init()

	m, err := migrate.New(migrationsFolder, settings.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected. Running migrations")

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	} else {
		log.Println("Database successfully migrated")
	}

	version, dirty, _ := m.Version()
	fmt.Println("version:", version)
	if dirty {
		fmt.Println("database state is currently DIRTY")
	}
}
