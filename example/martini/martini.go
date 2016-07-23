package main

import (
	"fmt"
	"github.com/codegangsta/martini"
)

func test(name string) {
	fmt.Println("ssss :", name)
}

func prefunc() {
	fmt.Println("aaaaaaaaaaa")
}

func main() {
	classic := martini.Classic()

	classic.Use(prefunc)
	classic.Map("szy")
	classic.Handle()
	classic.Invoke(test)
}
