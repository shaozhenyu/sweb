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

type RegisterArgs struct {
	Mobile   string `json:"mobile" valid:"required"`
	Password string `json:"password"`
	Code     string `json:"code" valid:"required"`
}

type ChatUserArgs struct {
	Name     string `json:"name" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func Register(install_ *install.Install, db *odm.DB) {

	install_.Group(ApiPrefix, func(r martini.Router) {
		r.Post("/get_verify_code", bind.Json2Struct(GetVerifyCodeArgs{}), GetVerifyCode)
		r.Post("/register", bind.Json2Struct(RegisterArgs{}), RegisterHandler)

		r.Post("/chat_register", bind.Json2Struct(ChatUserArgs{}), ChatRegisterHandler)
		r.Post("/chat_login", bind.Json2Struct(ChatUserArgs{}), ChatLoginHander)
	})

}
