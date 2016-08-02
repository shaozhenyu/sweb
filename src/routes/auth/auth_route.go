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

type User struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
	Password  string `json:"-" bson:"password" readonly:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func Register(install_ *install.Install, db *odm.DB) {

	install_.Group(ApiPrefix, func(r martini.Router) {
		r.Post("/get_verify_code", bind.Json2Struct(GetVerifyCodeArgs{}), GetVerifyCode)
		r.Post("/register", bind.Json2Struct(RegisterArgs{}), RegisterHandler)
	})

}
