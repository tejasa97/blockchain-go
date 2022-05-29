package node

import (
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
