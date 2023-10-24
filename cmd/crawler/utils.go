package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "modernc.org/sqlite"             // sqlite driver
)

func openDB(url string) (db *sql.DB, err error) {
	driver, conn := parseDbUrl(url)

	for attempts := 10; attempts > 0; attempts-- {
		db, err = sql.Open(driver, conn)
		if err != nil {
			time.Sleep(2 * time.Second)
			log.Warn("error opening", "db", url, "err", err)
		} else {
			break
		}
	}

	if err != nil {
		err = fmt.Errorf("error opening database %s: %w", url, err)
	}
	return
}

func parseDbUrl(url string) (driver, conn string) {
	ss := strings.SplitN(url, "://", 2)
	if len(ss) < 2 {
		driver, conn = "sqlite", url
	} else {
		driver, conn = ss[0], ss[1]
	}
	return
}
