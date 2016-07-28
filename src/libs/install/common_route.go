package install

import (
	"fmt"
	"net/http"
	"strconv"

	"libs/errorcode"
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
				return errorcode.HandleError(errorcode.ErrBadRequestBody)
			}

			statusCode, ret := GetResource(log, db, id, collname, req)
			return statusCode, ret
		})

		r.Post(fmt.Sprintf("/%s", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {

			statusCode, ret := PostResource(log, db, collname, req)
			return statusCode, ret
		})

		r.Get(fmt.Sprintf("/%s", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {

			m := (map[string]string)(params)
			statusCode, ret := ListResource(log, db, collname, m, req)
			return statusCode, ret
		})

		r.Put(fmt.Sprintf("/%s/(?P<id>[0-9]+$)", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {

			idStr := params["id"]
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return errorcode.HandleError(errorcode.ErrBadRequestBody)
			}

			m := (map[string]string)(params)

			statusCode, ret := PutResource(log, db, id, collname, req, m)
			return statusCode, ret
		})

		r.Delete(fmt.Sprintf("/%s/(?P<id>[0-9]+$)", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {
			idStr := params["id"]
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return errorcode.HandleError(errorcode.ErrBadRequestBody)
			}

			statusCode, ret := DeleteResource(log, db, id, collname, req)
			return statusCode, ret
		})
	})

}
