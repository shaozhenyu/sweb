package main

import (
	"os"

	"slsync/client"
	"slsync/server"

	"ember/cli"
)

func main() {
	cli := cli.NewCmds()

	client.Reg(cli)
	server.Reg(cli)

	args := os.Args[1:]
	if len(args) == 0 {
		cli.Help(true)
	} else {
		cli.Run(args)
	}
}
