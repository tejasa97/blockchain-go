package node

import (
	"context"
	"fmt"
	"github.com/tejasa97/go-block/database"
	"math/rand"
	"time"
)

const (
	miningDifficulty      = 1
	miningIntervalSeconds = 5
)

type PendingBlock struct {
	parent database.Hash
	number uint64
	time   uint64
	txs    []database.Tx
}

func NewPendingBlock(parentHash database.Hash, blockNumber uint64, txs []database.Tx) *PendingBlock {
	return &PendingBlock{
		parent: parentHash,
		number: blockNumber,
		time:   uint64(time.Now().Unix()),
		txs:    txs}
}

func createRandomPendingBlock() *PendingBlock {
	return NewPendingBlock(
		database.Hash{},
		0,
		[]database.Tx{
			database.NewTx("tejas", "rohit", 3, ""),
			database.NewTx("tejas", "tejas", 700, "reward"),
		})
}

func (n *Node) mine(ctx context.Context) error {
	var miningCtx context.Context
	var stopCurrentMining context.CancelFunc
	ticker := time.NewTicker(time.Second * miningIntervalSeconds)

	for {
		select {
		case <-ticker.C:
			go func() {
				fmt.Printf("checking if anything to mine ...")
				if len(n.pendingTXs) > 0 && !n.isMining {
					fmt.Printf("found stuff to mine! ...")
					n.isMining = true

					miningCtx, stopCurrentMining = context.WithCancel(ctx)
					err := n.minePendingTXs(miningCtx)
					if err != nil {
						fmt.Printf("ERROR: %s\n", err)
					}

					n.isMining = false
				}
			}()

		case block, _ := <-n.newSyncedBlocks:
			if n.isMining {
				blockHash, _ := block.Hash()
				fmt.Printf("\nPeer mined next Block `%s` faster :( \n", blockHash.Hex())

				n.removeMinedPendingTXs(block)
				stopCurrentMining()
			}

		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func Mine(ctx context.Context, pb PendingBlock) (database.Block, error) {

	if len(pb.txs) == 0 {
		err := fmt.Errorf("mining empty blocks is not allowed")
		return database.Block{}, err
	}

	// Init necessary variables
	start := time.Now()
	attempt := 0
	var block database.Block
	var hash database.Hash
	var nonce uint32

	// start a loop, repeating the hash generation until a valid block hash is found
	for !database.IsBlockHashValid(hash, miningDifficulty) {
		select {
		// in case the mining process needs to be stopped. eg: someone closes the program
		case <-ctx.Done():
			fmt.Println("Mining cancelled!")
			err := fmt.Errorf("mining cancelled. %s", ctx.Err())
			return database.Block{}, err
		default:
		}

		// generate a random big number
		attempt++
		nonce = generateNonce()

		if attempt%1000000 == 0 || attempt == 1 {
			fmt.Printf("Mining %d pending TXs, attempt: %d\n", len(pb.txs), attempt)
		}
		//	create a new block with this random nonce
		block = database.NewBlock(
			pb.parent,
			pb.number,
			pb.time,
			nonce,
			pb.txs)

		// hash and check?
		blockHash, err := block.Hash()
		if err != nil {
			err = fmt.Errorf("couldn't mine block, err: %s", err.Error())
			return database.Block{}, err
		}
		hash = blockHash
	}

	fmt.Printf("\nMined new Block '%x' using PoW with hash %x:\n", block, hash)
	fmt.Printf("\tHeight: '%v'\n", block.Header.Number)
	fmt.Printf("\tNonce: '%v'\n", block.Header.Nonce)
	fmt.Printf("\tCreated: '%v'\n", block.Header.Time)
	//fmt.Printf("\tMiner: '%v'\n", block.Header.Miner)
	fmt.Printf("\tParent: '%v'\n\n", block.Header.Parent.Hex())
	fmt.Printf("\tAttempt: '%v'\n", attempt)
	fmt.Printf("\tTime: %s\n\n", time.Since(start))

	return block, nil
}

func generateNonce() uint32 {
	rand.Seed(time.Now().UTC().UnixNano())

	return rand.Uint32()
}
