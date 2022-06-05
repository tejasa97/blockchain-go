package database

import (
	"crypto/sha256"
	"encoding/json"
)

type Account string

type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

func NewTx(from, to Account, value uint, data string) Tx {
	return Tx{from, to, value, data}
}

func (t Tx) IsReward() bool {
	return t.Data == "reward"
}

func (t Tx) Encode() ([]byte, error) {
	return json.Marshal(t)
}

func (t Tx) Hash() (Hash, error) {
	txJson, err := t.Encode()
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(txJson), nil
}

func NewAccount(value string) Account {
	return Account(value)
}
