package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	latestBlock     Block
	latestBlockHash Hash

	dbFile *os.File
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("Insufficient balance!")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) Copy() State {

	c := State{}
	c.latestBlock = s.latestBlock
	c.latestBlockHash = s.latestBlockHash
	c.Balances = make(map[Account]uint)

	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}

	return c
}
func (s *State) AddBlock(block Block) (Hash, error) {

	pendingState := s.Copy()

	err := pendingState.applyBlock(block)
	if err != nil {
		return Hash{}, err
	}

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{blockHash, block}
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("\nPersisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}

	s.Balances = pendingState.Balances
	s.latestBlock = pendingState.latestBlock
	s.latestBlockHash = blockHash

	return blockHash, nil
}

func (s *State) GetLatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) GetLatestBlockHeader() BlockHeader {
	return s.latestBlock.Header
}

func (s *State) Add(tx Tx) error {
	// Applies the transaction to the state
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Persist() (Hash, error) {

	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{blockHash, block}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash
	s.txMempool = []Tx{}

	return s.latestBlockHash, nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) applyBlock(b Block) error {
	// Applies the Block transactions to the state

	for _, tx := range b.TXs {
		if err := s.apply(tx); err != nil {
			return err
		}
	}

	return nil
}

// Gets the current State from Disk
func NewStateFromDisk() (*State, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genesis, err := loadGenesis()
	if err != nil {
		return nil, err
	}

	// Store the `balances` of each user
	balances := make(map[Account]uint)
	for account, balance := range genesis.Balances {
		balances[account] = balance
	}

	blockDbFilePath := filepath.Join(cwd, "database", "block.db")
	f, err := os.OpenFile(blockDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := NewState(balances, f)

	// Iterate over each the block.db file's line
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Deserialize JSON encoded Block
		var blockFs BlockFS
		err = json.Unmarshal(scanner.Bytes(), &blockFs)
		if err != nil {
			return nil, err
		}

		// Apply the block
		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		blockHash, err := blockFs.Value.Hash()
		if err != nil {
			return nil, err
		}

		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockHash
	}

	return state, nil
}
