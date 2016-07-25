package friends

import (
	"net/http"
	//"reflect"
	"strconv"

	"libs/odm"

	//"github.com/codegangsta/martini"
	"github.com/qiniu/xlog"
	"gopkg.in/mgo.v2/bson"
)

func GetFriends(log *xlog.Logger, db *odm.DB, r *http.Request) (int, interface{}) {
	log.Info("getfriends")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 400, nil
	}

	log.Info("id is : ", id)

	c := db.Session.DB("sweb").C("friends")

	var users []Friends
	err = c.Find(bson.M{"name": "szy"}).All(&users)
	if err != nil {
		return 400, nil
	}
	log.Info(users)
	return 200, users
}

func AddFriends(log *xlog.Logger, db *odm.DB, req *http.Request) (int, interface{}) {
	log.Info("add friends")

	v, err := db.Insert2("friends", req)
	if err != nil {
		return 400, err
	}

	return 200, v
}
