package node

import (
	"fmt"
	"net/http"
)

func (n *Node) serveHttp() error {
	// Init server, mux, and controller
	handler := http.NewServeMux()
	httpController := NewController()

	// URL mappings
	handler.HandleFunc(endpointNodeStatus, func(w http.ResponseWriter, r *http.Request) {
		httpController.getStatus(w, r, n)
	})
	handler.HandleFunc(endpointNodeBlocksQuery, func(w http.ResponseWriter, r *http.Request) {
		httpController.queryBlocks(w, r)
	})
	handler.HandleFunc(endpointNodeAdd, func(w http.ResponseWriter, r *http.Request) {
		httpController.addPeer(w, r, n)
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
