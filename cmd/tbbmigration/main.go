package main

import (
	"goblockchain/database"
	"os"
	"time"
)

func main() {
	state, err := database.NewStateFromDisk("")
	if err != nil {
		state.Logger.Error(err.Error())
		os.Exit(1)
	}
	defer state.Close()
	block0 := database.NewBlock(database.Hash{}, uint64(time.Now().Unix()), []database.Tx{
		database.NewTx(database.Account("rishi"), database.Account("rishi"), 700, ""),
	})
	err = state.AddBlock(block0)
	if err != nil {
		state.Logger.Error(err.Error())
		os.Exit(1)
	}
	block0Hash, err := state.Persist()
	if err != nil {
		state.Logger.Error(err.Error())
		os.Exit(1)
	}
	block1 := database.NewBlock(block0Hash, uint64(time.Now().Unix()), []database.Tx{
		database.NewTx("rishi", "baba", 2000, ""),
		database.NewTx("rishi", "rishi", 100, "reward"),
		database.NewTx("baba", "rishi", 1, ""),
		database.NewTx("baba", "caesar", 1000, ""),
		database.NewTx("baba", "rishi", 50, ""),
		database.NewTx("rishi", "rishi", 600, "reward"),
	})
	err = state.AddBlock(block1)
	if err != nil {
		state.Logger.Error(err.Error())
		os.Exit(1)
	}
	_, err = state.Persist()
	if err != nil {
		state.Logger.Error(err.Error())
		os.Exit(1)
	}
}
