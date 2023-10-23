package main

import (
	"time"

	"gopkg.in/urfave/cli.v1"
)

var (
	apiDBFlag = &cli.StringFlag{
		Name:  "api-db",
		Usage: "API SQLite file name",
	}
	apiListenAddrFlag = &cli.StringFlag{
		Name:  "addr",
		Usage: "Listening address",
		Value: "0.0.0.0:10000",
	}
	bootnodesFlag = &cli.StringSliceFlag{
		Name: "bootnodes",
		Usage: ("Comma separated nodes used for bootstrapping. " +
			"Defaults to hard-coded values for the selected network"),
	}
	crawlerDBFlag = &cli.StringFlag{
		Name:  "crawler-db",
		Usage: "Crawler SQLite file name",
	}
	dropNodesTimeFlag = &cli.DurationFlag{
		Name:  "drop-time",
		Usage: "Time to drop crawled nodes without any updates",
		Value: 24 * time.Hour,
	}
	geoipdbFlag = &cli.StringFlag{
		Name:  "geoipdb",
		Usage: "geoip2 database location",
	}
	listenAddrFlag = &cli.StringFlag{
		Name:  "addr",
		Usage: "Listening address",
		Value: "0.0.0.0:0",
	}
	nodedbFlag = &cli.StringFlag{
		Name:  "nodedb",
		Usage: "Nodes database location. Defaults to in memory database",
	}
	nodeFileFlag = &cli.StringFlag{
		Name:  "nodefile",
		Usage: "Path to a node file containing nodes to be crawled",
	}
	nodekeyFlag = &cli.StringFlag{
		Name:  "nodekey",
		Usage: "Hex-encoded node key",
	}
	nodeURLFlag = &cli.StringFlag{
		Name:  "nodeURL",
		Usage: "URL of the node you want to connect to",
		// Value: "http://localhost:8545",
	}
	timeoutFlag = &cli.DurationFlag{
		Name:  "timeout",
		Usage: "Timeout for the crawling in a round",
		Value: 5 * time.Minute,
	}
	workersFlag = &cli.Uint64Flag{
		Name:  "workers",
		Usage: "Number of workers to start for updating nodes",
		Value: 16,
	}
)
