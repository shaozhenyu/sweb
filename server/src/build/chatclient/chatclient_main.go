package main

import (
	"os"

	"chatclient/service"
	"ember/cli"
)

func main() {
	cmds := cli.NewCmds()
	Reg(cmds)

	args := os.Args[1:]
	if len(args) == 0 {
		cmds.Help(true)
	} else {
		cmds.Run(args)
	}
}

func Reg(cmds *cli.Cmds) {
	cmds.Reg("login", "user login", service.Login)
	cmds.Reg("register", "register new user", service.Register)

	login1 := cmds.Sub("login1", "user login another way")
	login1.Reg("chat", "chat with someone", service.ChatWith)
}
