package main

import (
	"fmt"
	"net/http"

	"models"
	"routes/uc"

	"libs/cache"
	"libs/install"
	"libs/odm"

	"github.com/go-martini/martini"
)

func main() {

	fmt.Println("start sweb")

	DBHost := "127.0.0.1:27017"
	db, err := odm.New(DBHost, "sweb", nil)
	if err != nil {
		return
	}

	ins := install.New()

	db.NewGroup(models.Friends{})

	RedisHost := "127.0.0.1:6379"
	rd_cache, err := cache.New(RedisHost, 0, 100)
	if err != nil {
		return
	}

	ins.Map(db)
	ins.Map(rd_cache)

	uc.Register(ins)
	ins.RegisterCommon(db, "friends", ins, Mhandle())
	ins.RunOnAddr(":8080")
}

func Mhandle() martini.Handler {
	fmt.Println("1111111111")

	return func(context martini.Context, req *http.Request) {
		fmt.Println("22222222")
		return
	}
}
