package database

import (
	"crypto/sha256"
	"encoding/json"
)

type Hash [32]byte
type BlockHeader struct {
	Parent Hash
	Time   uint64
}
type Block struct {
	Header BlockHeader
	Txs    []Tx
}
type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func NewBlock(parent Hash, time uint64, txs []Tx) Block {
	return Block{
		Header: BlockHeader{
			Parent: parent,
			Time:   time,
		},
		Txs: txs,
	}
}

func (b Block) Hash() (Hash, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(data), nil
}
