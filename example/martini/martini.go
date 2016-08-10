package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"log"
)

func main() {
	m := martini.Classic()

	m.Use(func(c martini.Context, log *log.Logger) {
		log.Println("before a request")

		c.Next()

		log.Println("after a request")
	})

	//simple get
	m.Get("/", func() (int, string) {
		fmt.Println("test")
		return 200, "hello martini"
	})

	//get route params
	m.Get("/hello_1/:name", func(params martini.Params) string {
		return "hello " + params["name"]
	})

	m.Get("/hello_2/**", func(params martini.Params) string {
		return "hello " + params["_1"]
	})

	//route group
	m.Group("/group", func(r martini.Router) {
		r.Get("/1", func() string {
			return "11"
		})
		r.Get("/2", func() string {
			return "22"
		})

	}, MiddleWare)

	//m.Run() //default addr 3000
	m.RunOnAddr(":8080")
}

func MiddleWare() {
	fmt.Println("I am middleware")
}

//map
