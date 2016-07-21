package main

import (
	"fmt"
	"github.com/codegangsta/martini"
)

func test(name string) {
	fmt.Println("ssss :", name)
}

func main() {
	classic := martini.Classic()

	classic.Map("szy")
	classic.Invoke(test)
}
