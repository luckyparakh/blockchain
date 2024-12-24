package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var tbb = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	tbb.AddCommand(balanceCmd())
	tbb.AddCommand(versionCmd)
	tbb.AddCommand(txCmd())
	err := tbb.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
