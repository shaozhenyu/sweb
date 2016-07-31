package auth

import (
	"libs/install"
	"libs/install/bind"
	"libs/odm"

	"github.com/go-martini/martini"
)

const (
	ApiPrefix = "/api/sweb"
)

type GetVerifyCodeArgs struct {
	Mobile string `json:"mobile" valid:"required"`
}

func Register(install_ *install.Install, db *odm.DB) {

	install_.Group(ApiPrefix, func(r martini.Router) {
		r.Post("/get_verify_code", bind.Json2Struct(GetVerifyCodeArgs{}), GetVerifyCode)
	})

}
