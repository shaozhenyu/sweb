package service

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func Chat() {
	for {
		read := bufio.NewReader(os.Stdin)
		buf, _, err := read.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		b := bytes.SplitN(buf, []byte(" "), 3)
		fmt.Println(b)
	}
}
