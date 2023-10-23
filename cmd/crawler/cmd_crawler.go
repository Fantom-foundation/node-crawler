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
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher"
	"github.com/Fantom-foundation/go-opera/integration/makefakegenesis"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/Fantom-foundation/go-opera/utils"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/oschwald/geoip2-golang"
	"gopkg.in/urfave/cli.v1"

	"github.com/ethereum/node-crawler/pkg/common"
	"github.com/ethereum/node-crawler/pkg/crawlerdb"
	"github.com/ethereum/node-crawler/pkg/p2p"
)

var (
	crawlerCommand = cli.Command{
		Name:   "crawl",
		Usage:  "Crawl the ethereum network",
		Action: crawlNodes,
		Flags: []cli.Flag{
			bootnodesFlag,
			crawlerDBFlag,
			geoipdbFlag,
			listenAddrFlag,
			nodeFileFlag,
			nodeURLFlag,
			nodedbFlag,
			nodekeyFlag,
			timeoutFlag,
			workersFlag,
			launcher.FakeNetFlag,
			launcher.GenesisFlag,
			launcher.ExperimentalGenesisFlag,
		},
	}
)

func crawlNodes(ctx *cli.Context) error {
	var inputSet common.NodeSet
	nodesFile := ctx.String(nodeFileFlag.Name)
	if nodesFile != "" && gethCommon.FileExist(nodesFile) {
		inputSet = common.LoadNodesJSON(nodesFile)
	}

	var db *sql.DB
	if ctx.IsSet(crawlerDBFlag.Name) {
		name := ctx.String(crawlerDBFlag.Name)
		db, err := openDB(name)
		if err != nil {
			panic(err)
		}
		log.Info("Connected to db")
		err = crawlerdb.InitDB(db)
		if err != nil {
			panic(err)
		}
	}

	var geoipDB *geoip2.Reader
	if geoipFile := ctx.String(geoipdbFlag.Name); geoipFile != "" {
		geoipDB, err := geoip2.Open(geoipFile)
		if err != nil {
			return err
		}
		defer geoipDB.Close()
	}

	var genesisStore *genesisstore.Store
	switch {
	case ctx.IsSet(launcher.FakeNetFlag.Name):
		_, num, err := parseFakeGen(ctx.String(launcher.FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", launcher.FakeNetFlag.Name, "err", err)
		}
		genesisStore = makefakegenesis.FakeGenesisStore(num, utils.ToFtm(1000000000), utils.ToFtm(5000000))
		defer genesisStore.Close()

	case ctx.IsSet(launcher.GenesisFlag.Name):
		genesisPath := ctx.String(launcher.GenesisFlag.Name)

		f, err := os.Open(genesisPath)
		if err != nil {
			panic(fmt.Errorf("Failed to open genesis file: %v", err))
		}
		defer f.Close()

		var genesisHashes genesis.Hashes
		genesisStore, genesisHashes, err = genesisstore.OpenGenesisStore(f)
		if err != nil {
			panic(fmt.Errorf("Failed to read genesis file: %v", err))
		}
		defer genesisStore.Close()

		// check if it's a trusted preset
		{
			g := genesisStore.Genesis()
			gHeader := genesis.Header{
				GenesisID:   g.GenesisID,
				NetworkID:   g.NetworkID,
				NetworkName: g.NetworkName,
			}
			experimental := true
			for _, allowed := range launcher.AllowedOperaGenesis {
				if allowed.Hashes.Equal(genesisHashes) && allowed.Header.Equal(gHeader) {
					log.Info("Genesis file is a known preset", "name", allowed.Name)
					experimental = false
					break
				}
			}
			if experimental {
				if ctx.Bool(launcher.ExperimentalGenesisFlag.Name) {
					log.Warn("Genesis file doesn't refer to any trusted preset")
				} else {
					panic(fmt.Errorf("Genesis file doesn't refer to any trusted preset. Enable experimental genesis with --genesis.allowExperimental"))
				}
			}
		}

	default:
		panic("Genesis expected!")
	}

	var bootnodes []string
	bootnodes = ctx.StringSlice(bootnodesFlag.Name)
	if len(bootnodes) > 0 {
		bootnodes = strings.Split(bootnodes[0], ",") // NOTE: workaround as StringSlice does not work properly
	}
	if len(bootnodes) < 1 {
		// defaults
		bootnodes = launcher.Bootnodes[genesisStore.Header().NetworkName]
	}

	crawler := p2p.NewCrawler(
		genesisStore,
		ctx.String(nodeURLFlag.Name),
		ctx.String(listenAddrFlag.Name),
		ctx.String(nodekeyFlag.Name),
		bootnodes,
		ctx.Duration(timeoutFlag.Name),
		ctx.Uint64(workersFlag.Name),
		db,
		geoipDB,
		ctx.String(nodedbFlag.Name),
	)

	crawler.Start(inputSet, func(updatedSet common.NodeSet) {
		if nodesFile != "" {
			updatedSet.WriteNodesJSON(nodesFile)
		}
	})
	defer crawler.Stop()

	wait()
	return nil
}

func wait() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
}

func parseFakeGen(s string) (id idx.ValidatorID, num idx.Validator, err error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		err = fmt.Errorf("use %%d/%%d format")
		return
	}

	var u32 uint64
	u32, err = strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return
	}
	id = idx.ValidatorID(u32)

	u32, err = strconv.ParseUint(parts[1], 10, 32)
	num = idx.Validator(u32)
	if num < 0 || idx.Validator(id) > num {
		err = fmt.Errorf("key-num should be in range from 1 to validators (<key-num>/<validators>), or should be zero for non-validator node")
		return
	}

	return
}
