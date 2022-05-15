package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const flagDataDir = "datadir"

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(txsCmd())
	tbbCmd.AddCommand(runCmd)

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired(flagDataDir)
}

//func getDataDirFromCmd(cmd *cobra.Command) string {
//	dataDir, _ := cmd.Flags().GetString(flagDataDir)
//
//	return fs.ExpandPath(dataDir)
//}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
