package node

import (
	"context"
	"fmt"
	"time"
)

//type Sync interface {
//	performSync() error
//}
//
//type sync struct {
//	apiClient apiClient
//	node      *Node
//}

var nodeApiClient = NewApiClient()

func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Searching for new Peers and blocks..")
			n.doSync()

		case <-ctx.Done():
			ticker.Stop()
		}
	}
	return nil
}

func (n *Node) doSync() error {

	for _, peer := range n.knownPeers {
		// Same node
		if n.isSelf(peer) {
			continue
		}

		peerState, err := nodeApiClient.getPeerStatus(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Printf("removing Peer %s from `KnownPeers`", peer.tcpAddress())
			n.removePeer(peer)
			continue
		}
		_, err = nodeApiClient.addPeer(*n, peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		n.updateBlocksFromPeer(peerState, peer)
		n.updateKnownPeers(peerState)
	}

	return nil
}

func (n *Node) updateBlocksFromPeer(peerState StateRes, peerNode PeerNode) error {
	/*
		Updates blocks from a Peer
	*/

	localBlockNumber := n.state.GetLatestBlockHeader().Number
	if peerState.BlockNumber <= localBlockNumber {
		return nil
	}

	// Update the new blocks
	newBlocksCount := peerState.BlockNumber - localBlockNumber
	fmt.Printf("Found %d new blocks from Peer %s\n", newBlocksCount, peerNode.tcpAddress())
	blocks, err := nodeApiClient.queryBlocks(peerNode, n.state.GetLatestBlockHash())
	if err != nil {
		fmt.Printf("Error querying blocks from Peer %s\n", peerNode.tcpAddress())
		return err
	}

	for _, block := range blocks {
		n.state.AddBlock(block)
	}

	return nil
}

func (n *Node) updateKnownPeers(peerState StateRes) error {
	for _, peerNode := range peerState.KnownPeers {
		n.addPeer(peerNode)
	}
	return nil
}
