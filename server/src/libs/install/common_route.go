package install

import (
	"fmt"
	"net/http"
	"strconv"

	"libs/errorcode"
	"libs/odm"

	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/qiniu/xlog"
)

const (
	prefix = "/api/sweb"
)

func getIUser(ctx martini.Context) (odm.IUser, error) {
	i := ctx.Get(inject.InterfaceOf((*odm.IUser)(nil)))
	if !i.IsValid() {
		return nil, errorcode.ErrNeedLogin
	}
	if i.IsNil() {
		return nil, errorcode.ErrNeedLogin
	}
	return i.Interface().(odm.IUser), nil
}

func (this *Install) RegisterCommon(db *odm.DB, collname string, mhandlers ...martini.Handler) {

	if db == nil {
		return
	}
	this.Map(db)

	coll := db.Coll[collname]
	this.Group(prefix, func(r martini.Router) {
		if coll.IsAllowMethod("GET") {
			r.Get(fmt.Sprintf("/%s/(?P<id>[0-9]+$)", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {

				idStr := params["id"]
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					return errorcode.HandleError(errorcode.ErrBadRequestBody)
				}

				statusCode, ret := GetResource(log, db, id, collname, req)
				return statusCode, ret
			})
		}

		if coll.IsAllowMethod("POST") {
			r.Post(fmt.Sprintf("/%s", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request, mctx martini.Context) (int, interface{}) {

				statusCode, ret := PostResource(log, db, collname, req)
				return statusCode, ret
			})
		}

		if coll.IsAllowMethod("LIST") {
			r.Get(fmt.Sprintf("/%s", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {

				m := (map[string]string)(params)
				statusCode, ret := ListResource(log, db, collname, m, req)
				return statusCode, ret
			})
		}

		if coll.IsAllowMethod("PUT") {
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
		}

		if coll.IsAllowMethod("DELETE") {
			r.Delete(fmt.Sprintf("/%s/(?P<id>[0-9]+$)", collname), func(log *xlog.Logger, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {
				idStr := params["id"]
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					return errorcode.HandleError(errorcode.ErrBadRequestBody)
				}

				statusCode, ret := DeleteResource(log, db, id, collname, req)
				return statusCode, ret
			})
		}
	}, mhandlers...)

}
