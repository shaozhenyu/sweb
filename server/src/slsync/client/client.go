package client

import (
	"fmt"
	"log"
	"os"

	"slsync/share"

	"ember/cli"
)

var (
	list = "./.slsync/.list"
	dir  = ".slsync"

	data = "data"
)

func Reg(cmds *cli.Cmds) {
	cmds.Reg("init", "init slsync config", CmdInit)
	cmds.Reg("sync", "sync file or dir", CmdSync)
}

func CmdInit(args []string) {
	err := share.InitConfig(dir, list)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(data, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func CmdSync(args []string) {
	fmt.Println(args)
	//UpdateLocalList()
}
