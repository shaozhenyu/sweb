package main

import (
	"fmt"

	"models"
	"routes/uc"

	"libs/cache"
	"libs/install"
	"libs/odm"
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
	ins.RegisterCommon(db, "friends", ins)
	ins.RunOnAddr(":8080")
}
