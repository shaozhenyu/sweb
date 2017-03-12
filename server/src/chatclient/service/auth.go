package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func getNameAndPassword(cmd string) (string, string, string) {
	fmt.Println("Please input user name: ")
	read := bufio.NewReader(os.Stdin)
	username, err := read.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Please input password: ")
	password, err := read.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	if cmd == "login" {
		return username, password, ""
	}

	fmt.Println("Please confirm password: ")
	confirm, err := read.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return username, password, confirm
}

func Register(args []string) {

	username, password, confirm := getNameAndPassword("register")

	if password != confirm {
		log.Fatal("Error: two passwords are not the same")
	}

	m := map[string]interface{}{
		"name":     strings.TrimRight(username, "\n"),
		"password": strings.TrimRight(password, "\n"),
	}

	body, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://127.0.0.1:8080/api/sweb/chat_register", "application/json;charset=utf-8", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}
	fmt.Println("Register ok : ", string(body))
	fmt.Println("You have login the schat, please chat with friends(send someone msg or someone msg)")
	Chat()
}

func Login(args []string) {
	username, password, _ := getNameAndPassword("login")
	fmt.Println(username, password)
}
