package uc

import (
	"github.com/go-martini/martini"
)

func Register(m martini.Router) {
	m.Get("/password", GetPassword)
	m.Post("/set_password", SetPassword)
}
