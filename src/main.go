package main

import (
	"net/http"

	"github.com/codegangsta/martini"
)

func main() {
	m := martini.Classic()
	m.Get("/", hellosweb)
	m.Get("/password", getpassword)
	m.Run()
}

func getpassword(r *http.Request) (int, interface{}) {
	app := r.URL.Query().Get("app")
	if app == "" {
		return 400, "请求格式错误"
	}
	return 200, nil
}

func hellosweb() string {
	return "hello world!"
}
