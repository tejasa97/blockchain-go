package node

import (
	"fmt"
	"github.com/tejasa97/go-block/database"
	"net/http"
)

type Node struct {
	state *database.State
}

const (
	HTTP_PORT = 8000
)

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
		httpController.getStatus(w, r, n.state)
	})

	handler.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		httpController.listBalances(w, r, n.state)
	})

	handler.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		httpController.addTx(w, r, n.state)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", HTTP_PORT), Handler: handler}

	err := server.ListenAndServe()
	// This shouldn't be an error!
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func NewNode() *Node {
	return &Node{}
}
