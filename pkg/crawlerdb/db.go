package crawlerdb

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/oschwald/geoip2-golang"
	beacon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/codec"

	"github.com/ethereum/node-crawler/pkg/common"
)

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS connections (
		ID              VARCHAR(66) NOT NULL,
		Now             DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		ClientType      TEXT,
		PK              TEXT,
		SoftwareVersion TEXT,
		Capabilities    TEXT,
		NetworkID       INTEGER,
		ForkID          TEXT,
		Blockheight     TEXT,
		HeadHash        TEXT,
		IP              TEXT,
		Country         TEXT,
		City            TEXT,
		Coordinates     TEXT,
		FirstSeen       TEXT,
		LastSeen        TEXT,
		Seq             INTEGER,
		Score           INTEGER,
		ConnType        TEXT,
		PRIMARY KEY (ID, Now)
	);
	`)
	return err
}

// ETH2 is a SSZ encoded field.
type ETH2 []byte

func (v ETH2) ENRKey() string { return "eth2" }

func UpdateNodes(db *sql.DB, geoipDB *geoip2.Reader, nodes []common.NodeJSON) error {
	log.Info("Writing connections to db", "count", len(nodes))

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO connections(
			ID,
			ClientType,
			PK,
			SoftwareVersion,
			Capabilities,
			NetworkID,
			ForkID,
			Blockheight,
			HeadHash,
			IP,
			Country,
			City,
			Coordinates,
			FirstSeen,
			LastSeen,
			Seq,
			Score,
			ConnType
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range nodes {
		info := &common.ClientInfo{}
		if n.Info != nil {
			info = n.Info
		}

		if info.ClientType == "" && n.TooManyPeers {
			info.ClientType = "tmp"
		}
		connType := ""
		var portUDP enr.UDP
		if n.N.Load(&portUDP) == nil {
			connType = "UDP"
		}
		var portTCP enr.TCP
		if n.N.Load(&portTCP) == nil {
			connType = "TCP"
		}
		fid := fmt.Sprintf("Hash: %v, Next %v", info.ForkID.Hash, info.ForkID.Next)

		var eth2 ETH2
		if n.N.Load(&eth2) == nil {
			info.ClientType = "eth2"
			var dat beacon.Eth2Data
			if err := dat.Deserialize(codec.NewDecodingReader(bytes.NewReader(eth2), uint64(len(eth2)))); err == nil {
				fid = fmt.Sprintf("Hash: %v, Next %v", dat.ForkDigest, dat.NextForkEpoch)
			}
		}
		var caps string
		for _, c := range info.Capabilities {
			caps = fmt.Sprintf("%v, %v", caps, c.String())
		}
		var pk string
		if n.N.Pubkey() != nil {
			pk = fmt.Sprintf("X: %v, Y: %v", n.N.Pubkey().X.String(), n.N.Pubkey().Y.String())
		}

		var country, city, loc string
		if geoipDB != nil {
			// parse GeoIp info
			ipRecord, err := geoipDB.City(n.N.IP())
			if err != nil {
				return err
			}
			country, city, loc =
				ipRecord.Country.Names["en"],
				ipRecord.City.Names["en"],
				fmt.Sprintf("%v,%v", ipRecord.Location.Latitude, ipRecord.Location.Longitude)
		}

		_, err = stmt.Exec(
			n.N.ID().String(),
			info.ClientType,
			pk,
			info.SoftwareVersion,
			caps,
			info.NetworkID,
			fid,
			info.Blockheight,
			info.HeadHash.String(),
			n.N.IP().String(),
			country,
			city,
			loc,
			n.FirstResponse.String(),
			n.LastResponse.String(),
			n.Seq,
			n.Score,
			connType,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
