package node

import (
	"github.com/tejasa97/go-block/database"
	"net/http"
)

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxAddRes struct {
	Success bool          `json:"success"`
	Hash    database.Hash `json:"hash"`
}

func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.GetLatestBlockHash(), state.Balances})
}

func addTxHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
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
