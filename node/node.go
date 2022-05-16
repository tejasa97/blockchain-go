package node

import (
	"fmt"
	"github.com/tejasa97/go-block/database"
	"net/http"
)

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

type Node struct {
	state *database.State
	port  uint64

	knownPeers []PeerNode
}

func (n *Node) Run() error {
	state, err := database.NewStateFromDisk()
	if err != nil {
		return err
	}
	defer state.Close()
	n.state = state

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
	handler.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		httpController.listBalances(w, r, n.state)
	})
	handler.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		httpController.addTx(w, r, n.state)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", n.port), Handler: handler}

	err := server.ListenAndServe()
	// This shouldn't be an error!
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func NewNode(port uint64, bootstrap PeerNode) *Node {
	return &Node{port: port, knownPeers: []PeerNode{bootstrap}}
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, isActive bool) *PeerNode {
	return &PeerNode{IP: ip, Port: port, IsBootstrap: isBootstrap, IsActive: isActive}
}
