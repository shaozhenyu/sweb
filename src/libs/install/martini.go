package install

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/qiniu/xlog"
)

type Install struct {
	*martini.ClassicMartini
}

func New() *Install {
	classic := martini.Classic()
	classic.Use(xlogFunc)
	classic.Map(rewriteReturnHandler())
	return &Install{classic}
}

func xlogFunc(ctx martini.Context, w http.ResponseWriter, r *http.Request) {
	xl := xlog.New(w, r)
	ctx.Map(xl)
}

func rewriteReturnHandler() martini.ReturnHandler {
	return func(ctx martini.Context, vals []reflect.Value) {

		rv := ctx.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
		res := rv.Interface().(http.ResponseWriter)

		res.Header().Set("Content-Type", "application/json;charset=utf-8")

		var responseVal reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			responseVal = vals[1]
		} else if len(vals) > 0 {
			responseVal = vals[0]
		}

		if responseVal.Kind() == reflect.String {
			res.Write([]byte(responseVal.String()))
			return
		}

		if reflect.TypeOf(responseVal).Kind() == reflect.Ptr {
			responseVal = responseVal.Elem()
		}

		i := responseVal.Interface()
		if i == nil {
			return
		}

		bytes, err := json.Marshal(i)
		if err != nil {
			bytes = []byte(`{"message": "` + err.Error() + `"}`)
		}
		res.Write(bytes)
	}
}
