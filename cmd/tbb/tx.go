package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tejasa97/go-block/database"
)

func txsCmd() *cobra.Command {

	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txsCmd.AddCommand(txnAddCmd())
	return txsCmd

}

func txnAddCmd() *cobra.Command {
	var txnAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new TX to database",
		Run: func(cmd *cobra.Command, args []string) {

			// Get flags
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			value, _ := cmd.Flags().GetUint("value")
			data, _ := cmd.Flags().GetString("data")

			// Execute operation
			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			tx := database.NewTx(fromAcc, toAcc, value, data)

			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Println(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			err = state.Add(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			snapshot, err := state.Persist()
			if err != nil {
				fmt.Println(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Printf("Snapshot: %x", snapshot)
			fmt.Println("TX successfully added to the ledger!")
		},
	}

	txnAddCmd.Flags().String("from", "", "From what account to send tokens")
	txnAddCmd.MarkFlagRequired("from")
	txnAddCmd.Flags().String("to", "", "To what account to send tokens")
	txnAddCmd.MarkFlagRequired("to")
	txnAddCmd.Flags().Uint("value", 0, "How many tokens to send")
	txnAddCmd.MarkFlagRequired("value")
	txnAddCmd.Flags().String("data", "", "what kind of a transaction is it?")

	return txnAddCmd
}
