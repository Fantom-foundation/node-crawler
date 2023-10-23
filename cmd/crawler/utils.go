package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/log"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "modernc.org/sqlite"             // sqlite driver
)

func openDB(path string) (db *sql.DB, err error) {
	dbUrl, err := url.Parse(path)
	if err != nil {
		err = fmt.Errorf("error opening database %s: %w", path, err)
		return
	}

	for attempts := 10; attempts > 0; attempts-- {

		if dbUrl.IsAbs() && len(dbUrl.Scheme) > 0 {
			db, err = sql.Open(dbUrl.Scheme, path)
		} else {
			db, err = sql.Open("sqlite", path)
		}

		if err != nil {
			time.Sleep(2 * time.Second)
			log.Warn("error opening", "db", path, "err", err)
		} else {
			break
		}
	}

	if err != nil {
		err = fmt.Errorf("error opening database %s: %w", path, err)
	}
	return
}
