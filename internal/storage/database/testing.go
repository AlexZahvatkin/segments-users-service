package database

import (
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestDB(t *testing.T, databaseUrl string) *Queries {
	t.Helper()

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	m, err := migrate.New("file://../../sql/postgresql/schema", databaseUrl)
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Down(); err != nil {
		if err != migrate.ErrNoChange {
			t.Fatal(err)
		}
	}

	if err := m.Up(); err != nil {
		t.Fatal(err)
	}

	return New(db)
}
