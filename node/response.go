package node

import "github.com/tejasa97/go-block/database"

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

type TxAddRes struct {
	Success bool          `json:"success"`
	Hash    database.Hash `json:"hash"`
}

type StatusRes struct {
	BlockHash   database.Hash `json:"block_hash"`
	BlockNumber uint64        `json:"block_number"`
}
