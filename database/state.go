package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type State struct {
	Balances  map[Account]uint
	txMemPool []Tx

	dbFile *os.File
	Logger *slog.Logger

	latestBlockHash Hash
}

func NewStateFromDisk(dataDir string) (*State, error) {

	var state = State{}
	state.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Add func, line and file name
		Level:     slog.LevelDebug,
	}))
	state.Logger.Info("New State from Disk")

	err := initDataDirIfNotExists(dataDir)
	if err != nil {
		return nil, err
	}

	// Read Genesis file
	gf, err := loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		return nil, err
	}
	b := make(map[Account]uint)
	for account, balance := range gf.Balances {
		b[account] = balance
	}
	state.Balances = b
	state.Logger.Debug(fmt.Sprintf("GB %v", state.Balances))

	f, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	state.dbFile = f
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var tx Tx
		err := json.Unmarshal(scanner.Bytes(), &tx)
		state.Logger.Debug(fmt.Sprintf("tx %v", tx))
		if err != nil {
			return nil, err
		}

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}
	return &state, nil
}

func (s *State) Close() {
	s.dbFile.Close()
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}
	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("%s", "Insufficient Balance")
	}
	s.Balances[tx.To] += tx.Value
	s.Balances[tx.From] -= tx.Value
	return nil
}
func (s *State) AddBlock(b Block) error {
	for _, tx := range b.Txs {
		if err := s.Add(tx); err != nil {
			return err
		}
	}
	return nil
}
func (s *State) Add(tx Tx) error {
	err := s.apply(tx)
	if err != nil {
		return err
	}
	s.txMemPool = append(s.txMemPool, tx)
	return nil
}
func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) Persist() (Hash, error) {
	b := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMemPool)
	blockHash, err := b.Hash()
	if err != nil {
		return Hash{}, err
	}
	blockFS := BlockFS{blockHash, b}
	data, err := json.Marshal(blockFS)
	if err != nil {
		return Hash{}, err
	}
	s.Logger.Info(fmt.Sprintf("Persisting new Block to disk:%v\n", blockFS))

	_, err = s.dbFile.Write(append(data, '\n'))
	if err != nil {
		return Hash{}, err
	}
	s.latestBlockHash = blockHash
	s.txMemPool = []Tx{}
	return blockHash, nil
}
