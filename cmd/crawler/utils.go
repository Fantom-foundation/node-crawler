package main

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"  // postgresql driver
	_ "modernc.org/sqlite" // sqlite driver
)

func openDB(path string) (db *sql.DB, err error) {
	dbUrl, err := url.Parse(path)
	if err != nil {
		err = fmt.Errorf("error opening database %s: %w", path, err)
		return
	}

	if dbUrl.IsAbs() && len(dbUrl.Scheme) > 0 {
		db, err = sql.Open(dbUrl.Scheme, path)
	} else {
		db, err = sql.Open("sqlite", path)
	}

	if err != nil {
		err = fmt.Errorf("error opening database %s: %w", path, err)
	}
	return
}
