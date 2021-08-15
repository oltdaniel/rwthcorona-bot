package utils

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DATA_BASE = os.Getenv("DATA_DIR")

var DATABASE *sql.DB = func() *sql.DB {
	database, err := sql.Open("sqlite3", filepath.Join(DATA_BASE, "/corona.db?cache=shared&mode=memory"))
	if err != nil {
		log.Fatal(err)
	}
	err = ensureMigrations(database)
	if err != nil {
		log.Fatal(err)
	}
	return database
}()

var migrations []string = []string{
	"CREATE TABLE IF NOT EXISTS corona_data (tag DATE, plz INTEGER, label TEXT, altersgruppe TEXT, anzahlWoche INTEGER, rateWoche REAL, anteilWoche REAL, UNIQUE(tag, plz, altersgruppe))",
}

func ensureMigrations(d *sql.DB) error {
	for _, table := range migrations {
		st, err := d.Prepare(table)
		if err != nil {
			return err
		}
		_, err = st.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
