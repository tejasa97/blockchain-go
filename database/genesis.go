package database

import (
	"encoding/json"
)

var genesisJson = `{
	"genesis_time": "2019-03-18T00:00:00.000000000Z",
	"chain_id": "the-blockchain-bar-ledger",
	"balances": {
		"tejas": 1000000
	}
}`

type Genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis() (Genesis, error) {
	content := []byte(genesisJson)

	var loadedGenesis Genesis
	err := json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return Genesis{}, err
	}

	return loadedGenesis, nil
}
