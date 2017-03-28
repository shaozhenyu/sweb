package server

import (
	"flag"
	"fmt"
	"log"
	"net"

	"slsync/share"

	"ember/cli"
)

var (
	list = "/tmp/slsync/.list"
)

func Reg(cmds *cli.Cmds) {
	cmds.Reg("run", "run slsync", CmdRun)
}

func CmdRun(args []string) {
	var dir string
	var port string

	newFlag := flag.NewFlagSet("", flag.ExitOnError)
	newFlag.StringVar(&dir, "dir", "/tmp/slsync", "data storage diretory")
	newFlag.StringVar(&port, "port", "9394", "listen port")
	cli.ParseFlag(newFlag, args, "dir", "port")

	fmt.Println("init service")
	NewService(dir, port)
}

func NewService(dir, port string) {
	err := share.InitConfig(dir, list)
	if err != nil {
		log.Fatal(err)
	}

	port = fmt.Sprintf(":%s", port)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	Server(ln)
}

func Server(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			log.Fatal("network error ", err)
		}

		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	defer conn.Close()
	fmt.Println("handler...")
}
