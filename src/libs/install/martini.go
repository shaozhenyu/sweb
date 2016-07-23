package install

import (
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/qiniu/xlog"
)

type Install struct {
	*martini.ClassicMartini
}

func New() *Install {
	classic := martini.Classic()
	classic.Use(xlogFunc)
	return &Install{classic}
}

func xlogFunc(ctx martini.Context, w http.ResponseWriter, r *http.Request) {
	xl := xlog.New(w, r)
	ctx.Map(xl)
}
