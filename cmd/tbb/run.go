package main

import (
	"fmt"
	"goblockchain/node"
	"os"

	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	var run = &cobra.Command{
		Use:   "run",
		Short: "Run the blockchain node	server",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			err := node.Run(dataDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	addDefaultRequiredFlag(run)
	return run
}
