// NOT IN USE as replaced with Run.go

package main

import (
	"fmt"
	"goblockchain/database"
	"os"

	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

func txCmd() *cobra.Command {
	var txCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (add)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	txCmd.AddCommand(txAddCmd())
	return txCmd
}
func txAddCmd() *cobra.Command {
	var txAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add transaction",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)

			tx := database.NewTx(database.Account(from), database.Account(to), value, data)
			state, err := database.NewStateFromDisk("")
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			defer state.Close()
			err = state.Add(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			_, err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Println("TX successfully added to the ledger.")
		},
	}
	txAddCmd.Flags().StringP(flagFrom, "f", "", "from which account to send token")
	txAddCmd.MarkFlagRequired(flagFrom)

	txAddCmd.Flags().StringP(flagTo, "t", "", "to which account token to be send")
	txAddCmd.MarkFlagRequired(flagTo)
	txAddCmd.Flags().UintP(flagValue, "v", 0, "how many token to send")
	txAddCmd.MarkFlagRequired(flagValue)
	txAddCmd.Flags().StringP(flagData, "d", "", "possible value 'reward'")

	return txAddCmd
}
