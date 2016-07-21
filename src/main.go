package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"libs/cache"
	"routes/uc"
)

func main() {
	fmt.Println("start sweb")
	m := martini.Classic()

	RedisHost := "127.0.0.1:6379"
	rd_cache, err := cache.New(RedisHost, 0, 100)
	if err != nil {
		return
	}
	m.Map(rd_cache)
	uc.Register(m)
	m.Run()
}
