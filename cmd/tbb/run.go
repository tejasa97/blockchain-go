package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tejasa97/go-block/node"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launches a Node and its HTTP api",
	Run: func(cmd *cobra.Command, args []string) {

		// Get flags
		isBootstrap, _ := cmd.Flags().GetBool("is_bootstrap")

		// Perform operation
		nodePort := uint64(8001)
		bootstrapNode := &node.PeerNode{}

		if !isBootstrap {
			bootstrapNode = node.NewPeerNode("localhost", node.BOOTSTRAP_NODE_PORT, true, true)
		} else {
			nodePort = node.BOOTSTRAP_NODE_PORT
		}

		node := node.NewNode("localhost", nodePort, isBootstrap, *bootstrapNode)
		err := node.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.PersistentFlags().Bool("is_bootstrap", false, "Is bootstrap node?")
}
