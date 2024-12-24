package main

import (
	"fmt"
	"goblockchain/database"
	"os"

	"github.com/spf13/cobra"
)

func balanceCmd() *cobra.Command{
	var balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Balance Operations",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	balanceCmd.AddCommand(balanceList)
	return balanceCmd
}

var balanceList = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer state.Close()
		fmt.Println("Accounts balances:")
		fmt.Println("__________________")
		for account, balance := range state.Balances {
			fmt.Printf("%s: %d\n", account, balance)
		}
	},
}
