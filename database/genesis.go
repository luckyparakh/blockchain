package database

import (
	"encoding/json"
	"os"
)

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
