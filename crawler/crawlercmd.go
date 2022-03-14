// Copyright 2021 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"

	"gopkg.in/urfave/cli.v1"
)

var (
	crawlerCommand = cli.Command{
		Name:      "crawl",
		Usage:     "Crawl the ethereum network",
		ArgsUsage: "<nodefile>",
		Action:    crawlNodes,
		Flags: []cli.Flag{
			utils.MainnetFlag,
			utils.RopstenFlag,
			utils.RinkebyFlag,
			utils.GoerliFlag,
			utils.NetworkIdFlag,
			bootnodesFlag,
			nodeURLFlag,
			nodeFileFlag,
			timeoutFlag,
			tableNameFlag,
			listenAddrFlag,
			nodekeyFlag,
			nodedbFlag,
		},
	}
	bootnodesFlag = cli.StringFlag{
		Name:  "bootnodes",
		Usage: "Comma separated nodes used for bootstrapping",
	}
	nodeURLFlag = cli.StringFlag{
		Name:  "nodeURL",
		Usage: "URL of the node you want to connect to",
		// Value: "http://localhost:8545",
	}
	nodeFileFlag = cli.StringFlag{
		Name:  "nodefile",
		Usage: "Path to a node file containing nodes to be crawled",
	}
	timeoutFlag = cli.DurationFlag{
		Name:  "timeout",
		Usage: "Timeout for the crawling in a round",
		Value: 5 * time.Minute,
	}
	tableNameFlag = cli.StringFlag{
		Name:  "table",
		Usage: "Name of the sqlite table",
	}
	listenAddrFlag = cli.StringFlag{
		Name:  "addr",
		Usage: "Listening address",
	}
	nodekeyFlag = cli.StringFlag{
		Name:  "nodekey",
		Usage: "Hex-encoded node key",
	}
	nodedbFlag = cli.StringFlag{
		Name:  "nodedb",
		Usage: "Nodes database location",
	}
)

func crawlNodes(ctx *cli.Context) error {
	var inputSet nodeSet

	if nodesFile := ctx.String(nodeFileFlag.Name); nodesFile != "" && common.FileExist(nodesFile) {
		inputSet = loadNodesJSON(nodesFile)
	}

	var db *sql.DB
	if ctx.IsSet(tableNameFlag.Name) {
		name := ctx.String(tableNameFlag.Name)
		shouldInit := false
		if _, err := os.Stat(name); os.IsNotExist(err) {
			shouldInit = true
		}
		var err error
		if db, err = sql.Open("sqlite3", name); err != nil {
			panic(err)
		}
		log.Info("Connected to db")
		if shouldInit {
			log.Info("DB did not exist, init")
			if err := createDB(db); err != nil {
				panic(err)
			}
		}
	}

	timeout := ctx.Duration(timeoutFlag.Name)

	nodeDB, err := enode.OpenDB(ctx.String(nodedbFlag.Name))
	if err != nil {
		panic(err)
	}

	for {
		inputSet = crawlRound(ctx, inputSet, db, nodeDB, timeout)
		if nodesFile := ctx.String(nodeFileFlag.Name); nodesFile != "" && common.FileExist(nodesFile) {
			writeNodesJSON(nodesFile, inputSet)
		}
	}
}

func crawlRound(ctx *cli.Context, inputSet nodeSet, db *sql.DB, nodeDB *enode.DB, timeout time.Duration) nodeSet {
	output := make(nodeSet)

	v5 := discv5(ctx, nodeDB, inputSet, timeout)
	output.add(v5.nodes()...)
	log.Info("DiscV5", "nodes", len(v5.nodes()))

	v4 := discv4(ctx, nodeDB, inputSet, timeout)
	output.add(v4.nodes()...)
	log.Info("DiscV4", "nodes", len(v4.nodes()))

	var nodes []nodeJSON
	for _, node := range output {
		nodes = append(nodes, node)
	}

	// Write the node info to influx
	if db != nil {
		if err := updateNodes(db, nodes); err != nil {
			panic(err)
		}
	}
	return output
}

func discv5(ctx *cli.Context, db *enode.DB, inputSet nodeSet, timeout time.Duration) nodeSet {
	ln, config := makeDiscoveryConfig(ctx, db)

	socket := listen(ln, ctx.String(listenAddrFlag.Name))

	disc, err := discover.ListenV5(socket, ln, config)
	if err != nil {
		panic(err)
	}
	defer disc.Close()

	return runCrawler(ctx, disc, inputSet, timeout)
}

func discv4(ctx *cli.Context, db *enode.DB, inputSet nodeSet, timeout time.Duration) nodeSet {
	ln, config := makeDiscoveryConfig(ctx, db)

	socket := listen(ln, ctx.String(listenAddrFlag.Name))

	disc, err := discover.ListenV4(socket, ln, config)
	if err != nil {
		panic(err)
	}
	defer disc.Close()

	return runCrawler(ctx, disc, inputSet, timeout)
}

func runCrawler(ctx *cli.Context, disc resolver, inputSet nodeSet, timeout time.Duration) nodeSet {
	genesis := makeGenesis(ctx)
	if genesis == nil {
		genesis = core.DefaultGenesisBlock()
	}
	networkID := ctx.Uint64(utils.NetworkIdFlag.Name)
	nodeURL := ctx.String(nodeURLFlag.Name)

	// Crawl the DHT for some time
	c := newCrawler(genesis, networkID, nodeURL, inputSet, disc, disc.RandomNodes())
	c.revalidateInterval = 10 * time.Minute
	return c.run(timeout)
}

// makeGenesis is the pendant to utils.MakeGenesis
// with local flags instead of global flags.
func makeGenesis(ctx *cli.Context) *core.Genesis {
	switch {
	case ctx.Bool(utils.RopstenFlag.Name):
		return core.DefaultRopstenGenesisBlock()
	case ctx.Bool(utils.RinkebyFlag.Name):
		return core.DefaultRinkebyGenesisBlock()
	case ctx.Bool(utils.GoerliFlag.Name):
		return core.DefaultGoerliGenesisBlock()
	default:
		return core.DefaultGenesisBlock()
	}
}
