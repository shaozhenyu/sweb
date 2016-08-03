package main

import (
	"fmt"
	//"net/http"

	"bind"
	"models"
	"routes/auth"
	"routes/uc"
	"server/idincr"

	"libs/cache"
	"libs/install"
	"libs/odm"

	//"github.com/go-martini/martini"
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

	db.NewGroup(models.Friends{}, models.MobileIdentity{}, models.User{})

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

	ins.RegisterCommon(db, "friends", ins, bind.BindAuthUser())
	ins.RegisterCommon(db, "mobileidentities", ins)
	ins.RegisterCommon(db, "users", ins)
	ins.RunOnAddr(":8080")
}
