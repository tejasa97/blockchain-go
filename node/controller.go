package node

import (
	"errors"
	"github.com/tejasa97/go-block/database"
	"net/http"
	"time"
)

const (
	endpointSyncQueryKeyFromBlock = "fromHash"
)

type Controller interface {
	listBalances(w http.ResponseWriter, r *http.Request, state *database.State)
	getStatus(w http.ResponseWriter, r *http.Request, node *Node)
	addTx(w http.ResponseWriter, r *http.Request, state *database.State)
	syncBlocks(w http.ResponseWriter, r *http.Request)
}

type controller struct {
}

func (h controller) listBalances(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.GetLatestBlockHash(), state.Balances})
}

func (h controller) getStatus(w http.ResponseWriter, r *http.Request, node *Node) {
	writeRes(w, node.getState())
}

func (h controller) syncBlocks(w http.ResponseWriter, r *http.Request) {

	reqHash := r.URL.Query().Get(endpointSyncQueryKeyFromBlock)
	if reqHash == "" {
		writeErrRes(w, errors.New("Missing blockhash"))
		return
	}

	hash := database.Hash{}
	err := hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		writeErrRes(w, errors.New("Invalid hash format"))
		return
	}

	// Read new blocks from the DB
	blocks, err := database.GetBlocksAfter(hash)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, SyncRes{Blocks: blocks})
}

func (h controller) addTx(w http.ResponseWriter, r *http.Request, state *database.State) {
	var req TxAddReq

	if err := readReq(r, &req); err != nil {
		writeErrRes(w, err)
		return
	}

	tx := database.NewTx(database.NewAccount(req.From), database.NewAccount(req.To), req.Value, req.Data)
	block := database.NewBlock(state.GetLatestBlockHash(), state.GetLatestBlockHeader().Number+1, uint64(time.Now().Unix()), []database.Tx{tx})
	hash, err := state.AddBlock(block)
	if err != nil {
		writeErrRes(w, err)
	}

	writeRes(w, TxAddRes{Hash: hash})
}

func NewController() Controller {
	return controller{}
}
