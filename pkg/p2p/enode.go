package p2p

import (
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/node-crawler/pkg/common"
)

func (c *Crawler) updateNode(list common.NodeSet, n *enode.Node, err error) {
	node, ok := list[n.ID()]

	// Skip validation of recently-seen nodes.
	if ok && !node.TooManyPeers {
		return
	}

	node.LastCheck = time.Now().UTC().Truncate(time.Second)

	if err != nil {
		if node.Score == 0 {
			// Node doesn't implement EIP-868.
			log.Debug("Skipping node", "id", n.ID())
			return
		}
		node.Score /= 2
	} else {
		node.N = n
		node.Seq = n.Seq()
		node.Score++
		if node.FirstResponse.IsZero() {
			node.FirstResponse = node.LastCheck
		}
		node.LastResponse = node.LastCheck
	}

	// Store/update node in output set.
	if node.Score <= 0 {
		log.Info("Removing node", "id", n.ID())
		delete(list, n.ID())
	} else {
		log.Info("Updating node", "id", n.ID(), "seq", n.Seq(), "score", node.Score)
		list[n.ID()] = node
	}
}
