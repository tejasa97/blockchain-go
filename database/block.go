package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
