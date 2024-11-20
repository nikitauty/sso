package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var user, password, host, port, dbname, migrationsPath, migrationsTable string

	flag.StringVar(&user, "user", "", "database user")
	flag.StringVar(&password, "password", "", "database password")
	flag.StringVar(&host, "host", "", "database host")
	flag.StringVar(&port, "port", "", "database port")
	flag.StringVar(&dbname, "dbname", "", "database name")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "table to migrate")
	flag.Parse()

	if user == "" {
		panic("user is required")
	}
	if password == "" {
		panic("password is required")
	}
	if host == "" {
		panic("host is required")
	}
	if port == "" {
		panic("port is required")
	}
	if dbname == "" {
		panic("dbname is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s:%s@%s:%s/dbname=%s?x-migrations-table=%s", user, password, host, port, dbname, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
