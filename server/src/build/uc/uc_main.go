package main

import (
	"fmt"
	//"net/http"

	"uc/bind"
	"uc/models"
	"uc/routes/auth"
	"uc/routes/city"
	"uc/routes/uc"
	"uc/server/idincr"

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

	db.NewGroup(models.Friends{AllowMethod: "GET|POST|LIST|PUT|DELETE"},
		models.MobileIdentity{},
		models.User{AllowMethod: "GET"},
		city.City{})

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
	city.Register(ins, db)

	ins.RegisterCommon(db, "friends", bind.BindAuthUser())
	ins.RegisterCommon(db, "mobileidentities")
	ins.RegisterCommon(db, "users")
	ins.RunOnAddr(":8080")
}
