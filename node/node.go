package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tejasa97/go-block/database"
)

const (
	BOOTSTRAP_NODE_IP   = "localhost"
	BOOTSTRAP_NODE_PORT = 8000
)

type Node struct {
	info         PeerNode
	state        *database.State
	pendingState *database.State

	pendingTXs  map[string]database.Tx
	archivedTXs map[string]database.Tx

	newSyncedBlocks chan database.Block
	isMining        bool

	knownPeers map[string]PeerNode
}

func NewNode(ip string, port uint64, isBootstrap bool, bootstrap PeerNode) *Node {

	knownPeers := make(map[string]PeerNode)
	if !isBootstrap {
		knownPeers[bootstrap.tcpAddress()] = bootstrap
	}

	info := NewPeerNode(ip, port, isBootstrap, true)
	return &Node{
		pendingTXs:      make(map[string]database.Tx),
		archivedTXs:     make(map[string]database.Tx),
		newSyncedBlocks: make(chan database.Block),
		info:            *info,
		knownPeers:      knownPeers,
		isMining:        false,
	}
}

func (n *Node) getState() StateRes {
	/*
		Returns the Node's state
	*/
	return StateRes{BlockHash: n.state.GetLatestBlockHash(), BlockNumber: n.state.GetLatestBlockHeader().Number, KnownPeers: n.knownPeers}
}

func (n *Node) addPeer(peerNode PeerNode) {
	/*
		Adds a peer to `knownPeers`
	*/

	if n.isSelf(peerNode) {
		return
	}

	_, isKnownPeer := n.knownPeers[peerNode.tcpAddress()]
	if !isKnownPeer {
		fmt.Printf("Found new peer %s \n", peerNode.tcpAddress())
		n.knownPeers[peerNode.tcpAddress()] = peerNode
	}
}

func (n *Node) removePeer(peer PeerNode) {
	/*
		Deletes a peer from `KnownPeers`
	*/

	delete(n.knownPeers, peer.tcpAddress())
}

func (n *Node) isSelf(peer PeerNode) bool {
	return n.info.tcpAddress() == peer.tcpAddress()
}

func (n *Node) addBlock(block database.Block) error {
	_, err := n.state.AddBlock(block)
	if err != nil {
		return err
	}

	// Reset the pending state
	pendingState := n.state.Copy()
	n.pendingState = &pendingState

	return nil
}

func (n *Node) getPendingTXsAsArray() []database.Tx {
	pendingTxs := make([]database.Tx, len(n.pendingTXs))

	idx := 0
	for _, tx := range n.pendingTXs {
		pendingTxs[idx] = tx
		idx++
	}

	return pendingTxs
}

func (n *Node) minePendingTXs(ctx context.Context) error {
	blockToMine := NewPendingBlock(
		n.state.GetLatestBlockHash(),
		n.state.GetNextBlockNumber(),
		n.getPendingTXsAsArray(),
	)

	minedBlock, err := Mine(ctx, *blockToMine)
	if err != nil {
		return err
	}

	n.removeMinedPendingTXs(minedBlock)

	err = n.addBlock(minedBlock)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) removeMinedPendingTXs(block database.Block) error {
	if len(block.TXs) == 0 || len(n.pendingTXs) == 0 {
		return nil
	}

	for _, tx := range block.TXs {
		txHash, _ := tx.Hash()
		if _, exists := n.pendingTXs[txHash.Hex()]; exists {
			fmt.Printf("\t-archiving mined TX: %s\n", txHash.Hex())

			n.archivedTXs[txHash.Hex()] = tx
			delete(n.pendingTXs, txHash.Hex())
		}
	}

	return nil
}

func (n *Node) AddPendingTX(tx database.Tx, fromPeer PeerNode) error {
	txHash, err := tx.Hash()
	if err != nil {
		return err
	}

	txJson, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	_, isAlreadyPending := n.pendingTXs[txHash.Hex()]
	_, isAlreadyArchived := n.archivedTXs[txHash.Hex()]

	if !(isAlreadyPending || isAlreadyArchived) {
		fmt.Printf("Added pending TX %s from Peer %s\n", txJson, fromPeer.tcpAddress())
		n.pendingTXs[txHash.Hex()] = tx
	}

	return nil
}

func (n *Node) Run() error {
	/*
		Runs the node's server and background sync job
	*/

	ctx := context.Background()
	state, err := database.NewStateFromDisk()
	if err != nil {
		return err
	}
	defer state.Close()
	n.state = state

	// Run sync in a goroutine
	go n.sync(ctx)
	go n.mine(ctx)

	fmt.Println("Blockchain state:")
	fmt.Printf("	- hash: %s\n", n.state.GetLatestBlockHash().Hex())

	return n.serveHttp()
}
