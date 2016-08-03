package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	h := sha1.New()
	h.Write([]byte("hello sha1"))
	fmt.Printf("%x", h.Sum(nil))
}
