package uc

import (
	"fmt"
	"libs/cache"
	"net/http"
)

func GetPassword(cache_ cache.CacheStorage, r *http.Request) (int, interface{}) {
	app := r.URL.Query().Get("app")
	if app == "" {
		return 400, "请求格式错误"
	}

	ps, err := cache_.Get(app)
	if err != nil {
		return 400, err
	}

	fmt.Println("password is ", string(ps))

	return 200, string(ps)
}

func SetPassword(cache_ cache.CacheStorage, r *http.Request) (int, interface{}) {
	app := r.URL.Query().Get("app")
	if app == "" {
		return 400, "请求格式错误"
	}
	passwd := r.URL.Query().Get("passwd")
	if passwd == "" {
		return 400, "请求格式错误"
	}

	if err := cache_.Set(app, []byte(passwd)); err != nil {
		return 400, "error"
	}

	fmt.Printf("app : %s, passwd : %s", app, passwd)
	return 200, "set ok"
}
