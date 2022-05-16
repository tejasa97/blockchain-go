package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tejasa97/go-block/node"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launches the Node and its HTTP api",
	Run: func(cmd *cobra.Command, args []string) {

		bootstrapNode := node.NewPeerNode("localhost", 8000, true, true)
		node := node.NewNode(8001, *bootstrapNode)
		err := node.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
