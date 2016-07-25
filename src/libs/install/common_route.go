package install

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"libs/odm"

// 	"github.com/codegangsta/inject"
// 	"github.com/codegangsta/martini"
// 	"github.com/qiniu/xlog"
// )

// func (this *Install) RegisterCommon(db *odm.DB, collname string, m *martini.Router) {
// 	m.Get(fmt.Sprintf("/%s/(?P<id>[0-9]+$)", collname), func(log *xlog, db *odm.DB, params martini.Params, req *http.Request) (int, interface{}) {
// 		idStr := params["id"]
// 		id, err := strconv.ParseInt(idStr, 10, 64)
// 		if err != nil {
// 			return err
// 		}

// 	})
// }
