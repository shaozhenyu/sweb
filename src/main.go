package main

import (
	"fmt"
	"net/http"

	"models"
	"routes/auth"
	"routes/uc"
	"server/idincr"

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

	coll := db.Session.DB("sweb").C("id_counter")
	db.SetIDMaker(idincr.NewIntIDMaker(coll))
	coll.Insert(odm.M{"_id": "sweb.friends", "offset": 1})

	db.NewGroup(models.Friends{}, models.MobileIdentity{})

	RedisHost := "127.0.0.1:6379"
	rd_cache, err := cache.New(RedisHost, 0, 100)
	if err != nil {
		return
	}

	ins := install.New()

	ins.Map(db)
	ins.Map(rd_cache)

	uc.Register(ins)
	auth.Register(ins, db)

	ins.RegisterCommon(db, "friends", ins, Mhandle())
	ins.RegisterCommon(db, "mobileidentities", ins)
	ins.RunOnAddr(":8080")
}

func Mhandle() martini.Handler {

	return func(context martini.Context, req *http.Request) {
		//TODO
		fmt.Println("do martini handler")
		return
	}
}
