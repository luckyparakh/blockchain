package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMemPool []Tx

	dbFile *os.File
	logger *slog.Logger
}

func NewStateFromDisk() (*State, error) {

	var state = State{}
	state.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Add func, line and file name
		Level: slog.LevelDebug,
	}))
	state.logger.Info("New State from Disk")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Read Genesis file
	gfPath := filepath.Join(cwd, "database", "genesis.json")
	state.logger.Debug(gfPath)
	gf, err := loadGenesis(gfPath)
	if err != nil {
		return nil, err
	}
	b := make(map[Account]uint)
	for account, balance := range gf.Balances {
		b[account] = balance
	}
	state.Balances = b
	state.logger.Debug(fmt.Sprintf("GB %v", state.Balances))
	
	txFilePath := filepath.Join(cwd, "database", "tx.db")
	state.logger.Debug(fmt.Sprintf("GB %v", txFilePath))
	f, err := os.OpenFile(txFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	state.dbFile = f
	scanner := bufio.NewScanner(f)
	
	for scanner.Scan() {
		var tx Tx
		err := json.Unmarshal(scanner.Bytes(), &tx)
		state.logger.Debug(fmt.Sprintf("tx %v", tx))
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

func (s *State) Add(tx Tx) error {
	err := s.apply(tx)
	if err != nil {
		return err
	}
	s.txMemPool = append(s.txMemPool, tx)
	return nil
}

func (s *State) Persist() error {
	memPool := make([]Tx, len(s.txMemPool))
	// Copy txMemPool as its value will change in the loop
	copy(memPool, s.txMemPool)
	for _, v := range memPool {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = s.dbFile.Write(append(data, '\n'))
		if err != nil {
			return err
		}
		// Remove the TX written to a file from the mempool
		s.txMemPool = s.txMemPool[1:]
	}
	return nil
}
