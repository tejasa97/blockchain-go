package node

import (
	"github.com/tejasa97/go-block/database"
	"net/http"
)

type Controller interface {
	listBalances(w http.ResponseWriter, r *http.Request, state *database.State)
	getStatus(w http.ResponseWriter, r *http.Request, node *Node)
	addTx(w http.ResponseWriter, r *http.Request, state *database.State)
}

type controller struct {
}

func (h controller) listBalances(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.GetLatestBlockHash(), state.Balances})
}

func (h controller) getStatus(w http.ResponseWriter, r *http.Request, node *Node) {
	writeRes(w, StatusRes{BlockHash: node.state.GetLatestBlockHash(), BlockNumber: node.state.GetLatestBlockHeader().Number, KnownPeers: node.knownPeers})
}

func (h controller) addTx(w http.ResponseWriter, r *http.Request, state *database.State) {
	var req TxAddReq

	if err := readReq(r, &req); err != nil {
		writeErrRes(w, err)
		return
	}

	tx := database.NewTx(database.NewAccount(req.From), database.NewAccount(req.To), req.Value, req.Data)
	err := state.Add(tx)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	// Persist to disk
	hash, err := state.Persist()
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, TxAddRes{Success: true, Hash: hash})
}

func NewController() Controller {
	return controller{}
}
