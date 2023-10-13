package crawler

import (
	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/inter/iblockproc"
	"github.com/Fantom-foundation/go-opera/inter/ier"
	"github.com/Fantom-foundation/go-opera/opera"
	"github.com/Fantom-foundation/go-opera/opera/genesisstore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type OperaStatus struct {
	NetworkName string
	Progress    gossip.PeerProgress
	NodeInfo    *gossip.NodeInfo
	Opera       *opera.Rules
	Chain       *params.ChainConfig
}

// LoadGenesis like gossip/Store.ApplyGenesis()
func (o *OperaStatus) LoadGenesis(gStore *genesisstore.Store) {
	var (
		g       = gStore.Genesis()
		hh      []opera.UpgradeHeight
		firstBS *iblockproc.BlockState
		firstES *iblockproc.EpochState
		lastES  *iblockproc.EpochState
	)
	g.Epochs.ForEach(func(er ier.LlrIdxFullEpochRecord) bool {
		es, bs := er.EpochState, er.BlockState

		if es.Rules.NetworkID != g.NetworkID || es.Rules.Name != g.NetworkName {
			panic("network ID/name mismatch")
		}

		if lastES == nil || es.Rules.Upgrades != lastES.Rules.Upgrades {
			hh = append(hh,
				opera.UpgradeHeight{
					Upgrades: es.Rules.Upgrades,
					Height:   bs.LastBlock.Idx + 1,
				})
		}
		lastES = &es
		if firstES == nil {
			firstES = &es
		}
		if firstBS == nil {
			firstBS = &bs
		}

		return true
	})

	if firstES == nil || firstBS == nil {
		panic("no ERs in genesis")
	}

	o.NetworkName = gStore.Header().NetworkName
	o.Progress = gossip.PeerProgress{
		Epoch:            firstES.Epoch,
		LastBlockIdx:     firstBS.LastBlock.Idx,
		LastBlockAtropos: firstBS.LastBlock.Atropos,
	}
	o.NodeInfo = &gossip.NodeInfo{
		Network:     g.NetworkID,
		Genesis:     common.Hash(g.GenesisID),
		Epoch:       firstES.Epoch,
		NumOfBlocks: firstBS.LastBlock.Idx,
	}
	o.Opera = &firstES.Rules
	o.Chain = firstES.Rules.EvmChainConfig(hh)
}
