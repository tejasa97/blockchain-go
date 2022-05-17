package node

import (
	"context"
	"fmt"
	"github.com/tejasa97/go-block/database"
	"net/http"
)

const (
	endpointNodeStatus = "/node/status"
	endpointSync       = "/node/sync"
	BOOTSTRAP_NODE_IP  = "localhost"
)
const (
	BOOTSTRAP_NODE_PORT = 8000
)

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

func (pn *PeerNode) tcpAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}

func (pn *PeerNode) apiProtocol() string {
	// Hardcode `apiProtocol` to "http" for now
	return "http"
}

type Node struct {
	info  PeerNode
	state *database.State

	knownPeers map[string]PeerNode
}

func (n *Node) Run() error {

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

func (n *Node) serveHttp() error {
	// Init server, mux, and controller
	handler := http.NewServeMux()
	httpController := NewController()

	// URL mappings
	handler.HandleFunc("/node/status", func(w http.ResponseWriter, r *http.Request) {
		httpController.getStatus(w, r, n)
	})
	handler.HandleFunc(endpointSync, func(w http.ResponseWriter, r *http.Request) {
		httpController.syncBlocks(w, r)
	})
	handler.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		httpController.listBalances(w, r, n.state)
	})
	handler.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		httpController.addTx(w, r, n.state)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", n.info.Port), Handler: handler}
	fmt.Println(fmt.Sprintf("Listening on port: %d", n.info.Port))

	err := server.ListenAndServe()
	// This shouldn't be an error!
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (n *Node) getState() StateRes {
	// Returns the Node's state
	return StateRes{BlockHash: n.state.GetLatestBlockHash(), BlockNumber: n.state.GetLatestBlockHeader().Number, KnownPeers: n.knownPeers}
}

func NewNode(ip string, port uint64, isBootstrap bool, bootstrap PeerNode) *Node {

	knownPeers := make(map[string]PeerNode)
	if !isBootstrap {
		knownPeers[bootstrap.tcpAddress()] = bootstrap
	}

	info := PeerNode{IP: ip, Port: port, IsBootstrap: isBootstrap, IsActive: true}
	return &Node{info: info, knownPeers: knownPeers}
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, isActive bool) *PeerNode {
	return &PeerNode{IP: ip, Port: port, IsBootstrap: isBootstrap, IsActive: isActive}
}

//func NewBootstrapNode() *Node {
//	knownPeers := make(map[string]PeerNode)
//	info := PeerNode{IP: BOOTSTRAP_NODE_IP, Port: BOOTSTRAP_NODE_PORT, IsBootstrap: true, IsActive: true}
//
//	return &Node{info: info, knownPeers: knownPeers}
//}
