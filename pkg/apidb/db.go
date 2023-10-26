package apidb

import (
	"database/sql"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/node-crawler/pkg/crawlerdb"
	"github.com/ethereum/node-crawler/pkg/vparser"
)

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS nodes (
			ID                  VARCHAR(66) NOT NULL,
			name                TEXT,
			version_major       INTEGER,
			version_minor       INTEGER,
			version_patch       INTEGER,
			version_tag         TEXT,
			version_build       TEXT,
			version_date        TEXT,
			os_name             TEXT,
			os_architecture     TEXT,
			language_name       TEXT,
			language_version    TEXT,
			last_crawled        DATETIME,
			country_name        TEXT,

			PRIMARY KEY (ID)
		);;
	`)
	return err
}

func InsertCrawledNodes(db *sql.DB, crawledNodes []crawlerdb.CrawledNode) error {
	log.Info("Writing nodes to db", "count", len(crawledNodes))

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO nodes(
			id,
			name,
			version_major,
			version_minor,
			version_patch,
			version_tag,
			version_build,
			version_date,
			os_name,
			os_architecture,
			language_name,
			language_version,
			last_crawled,
			country_name
		)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?) AS new
		ON DUPLICATE KEY UPDATE
			name = new.name,
			version_major = new.version_major,
			version_minor = new.version_minor,
			version_patch = new.version_patch,
			version_tag = new.version_tag,
			version_build = new.version_build,
			version_date = new.version_date,
			os_name = new.os_name,
			os_architecture = new.os_architecture,
			language_name = new.language_name,
			language_version = new.language_version,
			last_crawled = new.last_crawled,
			country_name = new.country_name
	`)
	if err != nil {
		return err
	}

	// It's possible for us to have the same node scraped multiple times, so
	// we want to make sure when we are upserting, we get the most recent
	// scrape upserted last.
	sort.SliceStable(crawledNodes, func(i, j int) bool {
		return crawledNodes[i].Now.Before(crawledNodes[j].Now)
	})

	for _, node := range crawledNodes {
		parsed := vparser.ParseVersionString(node.ClientType)
		if parsed != nil {
			_, err = stmt.Exec(
				node.ID,
				parsed.Name,
				parsed.Version.Major,
				parsed.Version.Minor,
				parsed.Version.Patch,
				parsed.Version.Tag,
				parsed.Version.Build,
				parsed.Version.Date,
				parsed.Os.Os,
				parsed.Os.Architecture,
				parsed.Language.Name,
				parsed.Language.Version,
				node.Now,
				node.Country,
			)
			if err != nil {
				panic(err)
			}
		} else {
			log.Warn("cann't parse client", "val", node.ClientType)
		}
	}
	return tx.Commit()
}

func DropOldNodes(db *sql.DB, minTimePassed time.Duration) error {
	log.Info("Dropping nodes", "older than", minTimePassed)
	oldest := time.Now().Add(-minTimePassed)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`DELETE FROM nodes WHERE last_crawled < ?`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(oldest)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	log.Info("Nodes drop", "affected", affected)
	return tx.Commit()
}
