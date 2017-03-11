package city

import (
	userBind "uc/bind"

	"libs/install"
	"libs/install/bind"
	"libs/odm"

	"github.com/go-martini/martini"
)

const (
	ApiPrefix = "/api/sweb"
)

type City struct {
	Id        int64  `json:"id" bson:"_id"`
	Uid       int64  `json:"uid" bson:"uid"`
	Name      string `json:"name" bson:"name" valid:"required"`
	Nation    string `json:"nation" bson:"nation" valid:"required"`
	Times     int64  `json:"times" bson:"times"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

func Register(install_ *install.Install, db *odm.DB) {
	install_.Group(ApiPrefix, func(r martini.Router) {
		r.Post("/new_city", bind.Json2Struct(City{}), AddNewCity)
		r.Get("/get_city_times", GetCityTimes)
	}, userBind.BindAuthUser())
}
