package node

import (
	"errors"
	"fmt"
	"github.com/tejasa97/go-block/database"
	"net/http"
)

const (
	endpointSyncQueryKeyFromBlock = "fromHash"
)

type Controller interface {
	listBalances(w http.ResponseWriter, r *http.Request, state *database.State)
	getStatus(w http.ResponseWriter, r *http.Request, node *Node)
	addTx(w http.ResponseWriter, r *http.Request, node *Node)
	queryBlocks(w http.ResponseWriter, r *http.Request)
	addPeer(w http.ResponseWriter, r *http.Request, node *Node)
}

type controller struct {
}

func (h controller) listBalances(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.GetLatestBlockHash(), state.Balances})
}

func (h controller) getStatus(w http.ResponseWriter, r *http.Request, node *Node) {
	writeRes(w, node.getState())
}

func (h controller) addPeer(w http.ResponseWriter, r *http.Request, node *Node) {
	var req AddNodeReq

	if err := readReq(r, &req); err != nil {
		writeErrRes(w, err)
		return
	}

	// Add the peer to the Node's `KnownPeers`
	req.IsActive = true
	node.addPeer(req.getPeerNode())

	writeRes(w, node.getState())
}

func (h controller) queryBlocks(w http.ResponseWriter, r *http.Request) {

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

	writeRes(w, BlocksRes{Blocks: blocks})
}

func (h controller) addTx(w http.ResponseWriter, r *http.Request, node *Node) {
	var req TxAddReq

	if err := readReq(r, &req); err != nil {
		writeErrRes(w, err)
		return
	}

	from := database.NewAccount(req.From)
	if from == "" {
		writeErrRes(w, fmt.Errorf("%s is an invalid 'from' sender", from))
		return
	}

	tx := database.NewTx(database.NewAccount(req.From), database.NewAccount(req.To), req.Value, req.Data)
	txHash, err := tx.Hash()
	if err != nil {
		writeErrRes(w, fmt.Errorf("Invalid TX hash"))
		return
	}

	err = node.AddPendingTX(tx, node.info)
	if err != nil {
		writeErrRes(w, err)
	}

	writeRes(w, TxAddRes{Success: true, Hash: txHash})
}

func NewController() Controller {
	return controller{}
}
