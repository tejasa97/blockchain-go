package node

import (
	"context"
	"encoding/hex"
	"github.com/tejasa97/go-block/database"
	"testing"
)

func TestValidBlockHash(t *testing.T) {
	hexHash := "00000000fa04f8160395c387277f8b2f14837603383d33809a4db586086168ed"

	var hash = database.Hash{}
	hex.Decode(hash[:], []byte(hexHash))

	isValid := database.IsBlockHashValid(hash, 4)
	if !isValid {
		t.Fatalf("hash '%s' with 6 zeroes should be valid", hexHash)
	}
}

func TestMine(t *testing.T) {
	pendingBlock := createRandomPendingBlock()

	//init context to be used to carry the deadline; i.e: get and act upon cancellation signals
	//and other request-scoped values between processes

	ctx := context.Background()

	minedBlock, err := Mine(ctx, *pendingBlock)
	if err != nil {
		t.Fatal(err)
	}

	minedBlockHash, err := minedBlock.Hash()
	if err != nil {
		t.Fatal(err)
	}

	if !database.IsBlockHashValid(minedBlockHash, miningDifficulty) {
		t.Fatal()
	}
}
