package database

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

func GetBlocksAfter(blockHash Hash) ([]Block, error) {
	// Returns all blocks after block with hash `blockHash`

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	blockDbFilePath := filepath.Join(cwd, "database", "block.db")
	f, err := os.OpenFile(blockDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	shouldStartCollecting := false
	blocks := make([]Block, 0)

	// Iterate over each the block.db file's line
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Unmarshal each block into `blockFs`
		var blockFs BlockFS
		err = json.Unmarshal(scanner.Bytes(), &blockFs)
		if err != nil {
			return nil, err
		}

		if shouldStartCollecting {
			blocks = append(blocks, blockFs.Value)
			continue
		}

		if blockHash == blockFs.Key {
			shouldStartCollecting = true
		}
	}
	return blocks, nil
}
