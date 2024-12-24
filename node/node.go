package node

import (
	"encoding/json"
	"fmt"
	"goblockchain/database"
	"net/http"
)

type BalanceRes struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balances"`
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

const httpPort = 8080

func Run(dataDir string) error {
	state, err := database.NewStateFromDisk(dataDir)
	if err != nil {
		return err
	}

	defer state.Close()
	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalances(w, r, state)
	})
	http.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, state)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func txAddHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	req := TxAddReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponseErr(w, err)
		return
	}
	tx := database.NewTx(database.Account(req.From), database.Account(req.To), req.Value, req.Data)
	err = state.Add(tx)
	if err != nil {
		writeResponseErr(w, err)
		return
	}
	hash, err := state.Persist()
	if err != nil {
		writeResponseErr(w, err)
		return
	}
	writeResponse(w, hash)
}

func listBalances(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeResponse(w, BalanceRes{state.LatestBlockHash(), state.Balances})
}

func writeResponse(w http.ResponseWriter, content any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(content)
}

func writeResponseErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
