package install

import (
	"fmt"
	"net/http"
	"strconv"

	"libs/odm"

	"github.com/go-martini/martini"
	"github.com/qiniu/xlog"
)

const (
	prefix = "/api/sweb"
)

func (this *Install) RegisterCommon(db *odm.DB, collname string, m martini.Router) {

	this.Group(prefix, func(r martini.Router) {
		r.Get(fmt.Sprintf("/%s/(?P<id>[0-9]+$)", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {

			idStr := params["id"]
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return 0, err
			}

			statusCode, ret := GetResource(log, db, id, collname, req)
			return statusCode, ret
		})

		r.Post(fmt.Sprintf("/%s", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {

			statusCode, ret := PostResource(log, db, collname, req)
			return statusCode, ret
		})
	})

}
