package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

type User struct {
	Name string
	Ip   string
	Conn net.Conn
}

var (
	allUser = map[string]*User{}
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		//每次连接的时候发一个用户信息
		buf := make([]byte, 1024)
		nameLen, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			log.Fatal(err)
		}

		username := string(buf[:nameLen])

		if len(username) == 0 {
			conn.Close()
			log.Fatal("get username error")
		}

		if _, ok := allUser[username]; ok {
			conn.Write([]byte("Error: user has login, can not login twice"))
			conn.Close()
			continue
		}

		newUser := &User{
			Name: username,
			Ip:   conn.RemoteAddr().String(),
			Conn: conn,
		}

		fmt.Println(username, ":", newUser.Ip, "connect")
		conn.Write([]byte("login ok"))

		allUser[username] = newUser

		go handleConnection(newUser)
	}
}

func handleConnection(user *User) {
	defer user.Conn.Close()
	sendMsgPrefix := user.Name + ": "
	for {
		buf := make([]byte, 1024)
		_, err := user.Conn.Read(buf)
		if err != nil {
			break
		}
		tmp := bytes.TrimLeft(buf, "send ")
		if len(tmp) < 2 {
			user.Conn.Write([]byte("please write the friend you need to chat"))
			continue
		}
		b := bytes.SplitN(tmp, []byte(" "), 2)
		friendName := string(b[0])
		if friendName == user.Name {
			user.Conn.Write([]byte("can not send msg to youself\n"))
			continue
		}
		if friend, ok := allUser[friendName]; ok {
			friend.Conn.Write([]byte(sendMsgPrefix + string(b[1])))
		} else {
			user.Conn.Write([]byte("friend not exists or offline\n"))
		}
	}
	delete(allUser, user.Name)
	fmt.Println(user.Name, ":", user.Ip, "disconnect")
}
