package service

import (
	"bufio"
	//"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

var ch chan int = make(chan int)

func ChatWith(args []string) {
	fmt.Println(args)
	println("TODO")
}

func Chat(username string) {
	conn, err := net.Dial("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(username))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	if string(buf[:5]) == "Error" {
		log.Fatal(string(buf))
	} else {
		fmt.Println(string(buf))
		fmt.Println("you can chat with friends(send someone msg or someone msg)")
	}

	go recvMsg(conn)

	read := bufio.NewReader(os.Stdin)
	for {
		buf, err := read.ReadBytes(byte('\n'))
		if err != nil {
			log.Fatal(err)
		}
		//conn.Write(bytes.TrimRight(buf, "\n"))
		conn.Write(buf)
		select {
		case <-ch:
			log.Fatal("recvMsg error")
		default:
		}

	}
}

func recvMsg(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			ch <- 1
			log.Fatal(err)
		}
		fmt.Printf(string(buf[:n]))
	}
}
