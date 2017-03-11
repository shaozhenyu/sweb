package city

import (
	//"fmt"
	"net/http"

	"libs/errorcode"
	"libs/odm"

	"github.com/qiniu/xlog"
	"gopkg.in/mgo.v2"
)

//TODO 单词首字母需要大写， 比如BeiJing China
func AddNewCity(log *xlog.Logger, db *odm.DB, authUser odm.IUser, args *City, r *http.Request) (int, interface{}) {

	m := City{}

	coll := db.C(m)
	defer coll.Close()

	cityName := args.Name
	cityNation := args.Nation

	err := coll.Find(odm.M{"name": cityName, "nation": cityNation}).One(&m)
	if err == nil {
		m.Times += 1
		if err = coll.Update(odm.M{"name": cityName, "nation": cityNation}, odm.M{"$set": odm.M{"times": m.Times}}); err != nil {
			return errorcode.HandleError(err)
		}
		return 200, m
	}
	if err != nil && err != mgo.ErrNotFound {
		return errorcode.HandleError(err)
	}

	{
		m.Uid = authUser.Uid()
		m.Name = cityName
		m.Nation = cityNation
		m.Times = 1
	}

	if err1 := db.Insert(&m, nil); err1 != nil {
		return errorcode.HandleError(err1)
	}

	return 200, m
}

func GetCityTimes(log *xlog.Logger, db *odm.DB, authUser odm.IUser, r *http.Request) (int, interface{}) {

	cityName := r.URL.Query().Get("name")
	cityNation := r.URL.Query().Get("nation")

	if cityName == "" || cityNation == "" {
		return errorcode.HandleError(errorcode.ErrBadRequestBody)
	}

	m := City{}
	coll := db.C(m)
	defer coll.Close()

	if err := coll.Find(odm.M{"name": cityName, "nation": cityNation}).One(&m); err != nil {
		return errorcode.HandleError(err)
	}

	return 200, m.Times
}
