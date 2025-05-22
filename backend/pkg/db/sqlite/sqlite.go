package sqlite

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	sqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDB() *sql.DB {
	databasePath := "pkg/db/social-network.db"

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalln("Error opening DB:", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err := runMigrations(db); err != nil {
		log.Fatalln("Migration error:", err)
	}

	return db
}

func GetDB() *sql.DB {
	return db
}


func runMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/db/migrations/sqlite",
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
