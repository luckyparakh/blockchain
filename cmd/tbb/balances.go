package main

import (
	"fmt"
	"goblockchain/database"
	"os"

	"github.com/spf13/cobra"
)

func balanceCmd() *cobra.Command {
	var balanceCmd = &cobra.Command{
		Use:   "balances",
		Short: "Balance Operations",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	balanceCmd.AddCommand(balancesListCmd())
	return balanceCmd
}
func balancesListCmd() *cobra.Command {
	var balanceList = &cobra.Command{
		Use:   "list",
		Short: "Lists all balances",
		Run: func(cmd *cobra.Command, args []string) {
			datadir, _ := cmd.Flags().GetString(flagDataDir)
			state, err := database.NewStateFromDisk(datadir)
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
	addDefaultRequiredFlag(balanceList)
	return balanceList
}
