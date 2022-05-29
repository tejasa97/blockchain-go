package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Hash [32]byte

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.Hex()), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Number uint64 `json:"number"`
	Nonce  uint32 `json:"nonce"`
	Time   uint64 `json:"time"`
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func NewBlock(parentHash Hash, blockNumber, time uint64, txs []Tx) Block {
	return Block{BlockHeader{Parent: parentHash, Number: blockNumber, Time: time}, txs}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}

func IsBlockHashValid(hash Hash, miningDifficulty uint) bool {
	// First `miningDifficulty` chars of the Hash need to be "0"

	for i := uint(0); i < miningDifficulty; i++ {
		if fmt.Sprintf("%x", hash[i]) != "0" {
			return false
		}
	}

	if fmt.Sprintf("%x", hash[miningDifficulty]) == "0" {
		return false
	}

	return true
}
