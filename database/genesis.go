package database

import (
	"encoding/json"
	"os"
)

var genesisJson = `
{
    "genesis_time": "2024-12-18T00:00:00.000000000Z",
    "chain_id": "the-blockchain-bar-ledger",
    "balances": {
      "rishi": 1000000
    }
}
`

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (*genesis, error) {
	gfData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data genesis
	err = json.Unmarshal(gfData, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func writeGenesisToDisk(path string) error {
	return os.WriteFile(path, []byte(genesisJson), 0644)
}
