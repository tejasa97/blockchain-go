package database

import "os"

func NewState(balances map[Account]uint, file *os.File) *State {
	return &State{balances, make([]Tx, 0), Block{}, Hash{}, file}
}
