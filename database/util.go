package database

import "os"

func NewState(balances map[Account]uint, file *os.File) *State {
	/*
		Returns a new `State` object
	*/

	return &State{
		Balances:        balances,
		txMempool:       make([]Tx, 0),
		latestBlock:     Block{},
		latestBlockHash: Hash{},
		hasGenesisBlock: false,
		dbFile:          file,
	}
}
