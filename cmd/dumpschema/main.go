package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/contextwtf/lanyard/api/migrations"
	"github.com/contextwtf/lanyard/migrate"
	_ "github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const temporaryDatabase = `tmp_db_for_dump_schema`

func main() {
	check(run("dropdb", "--if-exists", temporaryDatabase))
	check(run("createdb", temporaryDatabase))

	db, err := sql.Open("postgres", fmt.Sprintf("postgres:///%s?sslmode=disable", temporaryDatabase))
	check(err)
	_, err = db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, "public"))
	check(err)
	_, err = db.Exec(fmt.Sprintf(`ALTER DATABASE %s SET TIMEZONE TO 'UTC'`, temporaryDatabase))
	check(err)
	db.Close()

	db, err = sql.Open("postgres", fmt.Sprintf("postgres:///%s?sslmode=disable&search_path=%s", temporaryDatabase, "public"))
	check(err)
	check(migrate.Run(context.Background(), db, migrations.Migrations))

	var buf bytes.Buffer
	pgdump := exec.Command("pg_dump", "-sOx", temporaryDatabase)
	pgdump.Stdout = &buf
	pgdump.Stderr = os.Stderr
	check(pgdump.Run())

	f, err := os.Create(filepath.Join("api", "schema.sql"))
	check(err)
	defer f.Close()

	for _, line := range strings.Split(buf.String(), "\n") {
		if strings.HasPrefix(line, "--") || strings.Contains(line, "COMMENT") {
			continue
		}
		_, err = f.WriteString(line + "\n")
		check(err)
	}
	check(db.Close())
	check(run("dropdb", temporaryDatabase))
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
