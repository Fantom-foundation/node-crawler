package crawlerdb

import (
	"database/sql"
	"fmt"
	"time"
)

type CrawledNode struct {
	ID              string
	Now             time.Time
	ClientType      string
	SoftwareVersion uint64
	Capabilities    string
	NetworkID       uint64
	Country         string
	ForkID          string
}

func ReadAndDeleteUnseenNodes(db *sql.Tx) ([]CrawledNode, error) {
	rows, err := db.Query(`
		SELECT
			ID,
			Now,
			ClientType,
			SoftwareVersion,
			Capabilities,
			NetworkID,
			Country,
			ForkID
		FROM connections
	`)
	if err != nil {
		return nil, err
	}

	var nodes []CrawledNode
	for rows.Next() {
		var node CrawledNode
		err = rows.Scan(
			&node.ID,
			&timeScanner{&node.Now},
			&node.ClientType,
			&node.SoftwareVersion,
			&node.Capabilities,
			&node.NetworkID,
			&node.Country,
			&node.ForkID,
		)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	_, err = db.Exec(`
		DELETE FROM connections
	`)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// timeScanner implements sql.Scanner interface
type timeScanner struct {
	*time.Time
}

func (t *timeScanner) Scan(src interface{}) (err error) {
	const layout = "2006-01-02 15:04:05"

	switch v := src.(type) {
	case time.Time:
		*t.Time = v
	case []uint8:
		*t.Time, err = time.Parse(layout, string(v))
	case string:
		*t.Time, err = time.Parse(layout, v)
	default:
		err = fmt.Errorf("error while scan time: got %T from db", v)
	}

	return
}
