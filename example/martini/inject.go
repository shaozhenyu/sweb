package main

import (
	"fmt"
	"github.com/codegangsta/inject"
)

type SpecialString interface{}

func Say(name string, gender SpecialString, age int) {
	fmt.Printf("My name is %s, gender is %s, age is %d!\n", name, gender, age)
}

func main() {
	inj := inject.New()
	inj.Map("szy")
	inj.MapTo("man", (*SpecialString)(nil))
	inj.Map(20)
	inj.Invoke(Say)
}
