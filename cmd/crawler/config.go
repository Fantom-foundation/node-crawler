package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher"
	"github.com/Fantom-foundation/go-opera/integration/makefakegenesis"
	"github.com/Fantom-foundation/go-opera/opera/genesis"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/Fantom-foundation/go-opera/utils"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	"github.com/ethereum/node-crawler/pkg/crawler"
)

func mayGetOperaStatus(ctx *cli.Context) *crawler.OperaStatus {
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
			for _, allowed := range launcher.AllowedOperaGenesis {
				if allowed.Hashes.Equal(genesisHashes) && allowed.Header.Equal(gHeader) {
					log.Info("Genesis file is a known preset", "name", allowed.Name)
					goto notExperimental
				}
			}
			if ctx.Bool(launcher.ExperimentalGenesisFlag.Name) {
				log.Warn("Genesis file doesn't refer to any trusted preset")
			} else {
				panic(fmt.Errorf("Genesis file doesn't refer to any trusted preset. Enable experimental genesis with --genesis.allowExperimental"))
			}
		notExperimental:
		}

	default:
		panic("Genesis expected!")
	}

	operaStatus := new(crawler.OperaStatus)
	operaStatus.LoadGenesis(genesisStore)
	return operaStatus
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
