package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "processor error: %s", err)
		debug.PrintStack()
		os.Exit(1)
	}
}

func main() {
	const (
		defaultPGURL = "postgres:///al"
		migDir       = "file://db/sql/migrations"
	)
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		dburl = defaultPGURL
	}

	m, err := migrate.New(migDir, dburl)
	check(err)

	fmt.Println("migrating database")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		check(err)
	} else {
		fmt.Println("success")
	}

	version, dirty, _ := m.Version()
	fmt.Println("version:", version)
	if dirty {
		fmt.Println("database is dirty")
	}
}
