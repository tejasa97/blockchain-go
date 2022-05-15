package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tejasa97/go-block/database"
)

func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Prepare the blocks
	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("tejas", "tejas", 3, ""),
			database.NewTx("tejas", "tejas", 700, "reward"),
		},
	)

	state.AddBlock(block0)
	block0Hash := state.GetLatestBlockHash()
	fmt.Printf("latest block hash %x", block0Hash)

	defer state.Close()

	block1 := database.NewBlock(
		block0Hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("tejas", "rohit", 2000, ""),
			database.NewTx("tejas", "tejas", 100, "reward"),
			database.NewTx("rohit", "tejas", 1, ""),
			database.NewTx("rohit", "caesar", 1000, ""),
			database.NewTx("rohit", "tejas", 50, ""),
			database.NewTx("tejas", "tejas", 600, "reward"),
		},
	)

	state.AddBlock(block1)
}
