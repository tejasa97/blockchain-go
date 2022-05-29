package node

import (
	"context"
	"fmt"
	"github.com/tejasa97/go-block/database"
)

const (
	BOOTSTRAP_NODE_IP   = "localhost"
	BOOTSTRAP_NODE_PORT = 8000
)

type Node struct {
	info  PeerNode
	state *database.State

	knownPeers map[string]PeerNode
}

func NewNode(ip string, port uint64, isBootstrap bool, bootstrap PeerNode) *Node {

	knownPeers := make(map[string]PeerNode)
	if !isBootstrap {
		knownPeers[bootstrap.tcpAddress()] = bootstrap
	}

	info := PeerNode{IP: ip, Port: port, IsBootstrap: isBootstrap, IsActive: true}
	return &Node{info: info, knownPeers: knownPeers}
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

	fmt.Println("Blockchain state:")
	fmt.Printf("	- hash: %s\n", n.state.GetLatestBlockHash().Hex())

	return n.serveHttp()
}
