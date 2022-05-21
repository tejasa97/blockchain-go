package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	latestBlock     Block
	latestBlockHash Hash
	hasGenesisBlock bool

	dbFile *os.File
}

func (s *State) Copy() State {
	/*
		Creates a deepcopy of the `State` object
	*/

	c := State{}
	c.latestBlock = s.latestBlock
	c.latestBlockHash = s.latestBlockHash
	c.txMempool = make([]Tx, len(s.txMempool))
	c.Balances = make(map[Account]uint)

	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}

	for _, tx := range s.txMempool {
		c.txMempool = append(c.txMempool, tx)
	}

	return c
}
func (s *State) AddBlock(block Block) (Hash, error) {
	/*
		Adds a `block` to the ledger (block.db)
	*/

	pendingState := s.Copy()
	err := applyBlock(block, pendingState)
	if err != nil {
		return Hash{}, err
	}

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{Key: blockHash, Value: block}
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
	s.latestBlock = block
	s.latestBlockHash = blockHash
	s.hasGenesisBlock = true

	return blockHash, nil
}

func (s *State) GetLatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) GetLatestBlockHeader() BlockHeader {
	return s.latestBlock.Header
}

func (s *State) Persist() (Hash, error) {

	block := NewBlock(s.latestBlockHash, s.GetLatestBlockHeader().Number+1, uint64(time.Now().Unix()), s.txMempool)

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
	s.latestBlock = block

	return s.latestBlockHash, nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func applyBlock(b Block, state State) error {
	// Applies the Block transactions to the state
	// Block metadata is verified, as well as the transactions within - sufficient balances, etc

	nextExpectedBlockNumber := state.GetLatestBlockHeader().Number + 1

	if state.hasGenesisBlock && b.Header.Number != nextExpectedBlockNumber {
		return fmt.Errorf(
			"next expected block number should be `$d` not `%d`",
			nextExpectedBlockNumber,
			b.Header.Number,
		)
	}

	// validate incoming block parent hash equals the current (latest known) hash
	if state.hasGenesisBlock && state.GetLatestBlockHeader().Number > 0 && !reflect.DeepEqual(b.Header.Parent, state.latestBlockHash) {
		return fmt.Errorf(
			"next block's parent hash must be `%x` not `%x`",
			state.latestBlockHash,
			b.Header.Parent,
		)
	}

	return applyTXs(b.TXs, &state)
}

func applyTXs(txs []Tx, s *State) error {
	// Applies all Transactions `txs` to the State `s`

	for _, tx := range txs {
		err := ApplyTx(tx, s)
		if err != nil {
			return err
		}
	}

	return nil
}

// Applies a transaction `tx` to a State `s`
func ApplyTx(tx Tx, s *State) error {

	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	// Validate sender balance
	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf(
			"invalid TX. Sender `%s`'s balance is less than TX value (`%d` < `%d`)",
			tx.From,
			s.Balances[tx.From],
			tx.Value,
		)
	}

	// Update balances
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

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
		err = applyBlock(blockFs.Value, *state)
		if err != nil {
			return nil, err
		}

		blockHash, err := blockFs.Value.Hash()
		if err != nil {
			return nil, err
		}

		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockHash
		state.hasGenesisBlock = true
	}

	return state, nil
}
