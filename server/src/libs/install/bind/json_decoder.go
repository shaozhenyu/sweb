package bind

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"libs/errorcode"

	validator "github.com/asaskevich/govalidator"
	"github.com/go-martini/martini"
	"github.com/qiniu/log"
)

var (
	bodyLimiter = 8 << 20
)

func Json2Struct(obj interface{}) martini.Handler {
	objType := reflect.TypeOf(obj)
	if objType.Kind() != reflect.Struct {
		panic("obj must a struct, not a ptr")
	}

	return func(context martini.Context, req *http.Request) {
		if req.Body != nil {
			defer req.Body.Close()

			struct_ := reflect.New(objType).Interface()
			if err := Map2Json(struct_, req.Body); err != nil {
				log.Error(err)
				return
			}

			context.Map(struct_)
			return
		}
	}
}

func Map2Json(obj interface{}, r io.Reader) error {
	limitReader := io.LimitReader(r, int64(bodyLimiter))

	decoder := json.NewDecoder(limitReader)
	if err := decoder.Decode(obj); err != nil {
		log.Error(err)
		if err == io.ErrUnexpectedEOF {
			return errorcode.ErrTooLargeRequestBody
		}
		return errorcode.ErrBadRequestBody
	}

	if ok, err := validator.ValidateStruct(obj); !ok {
		log.Error(err)
		return errorcode.ErrBadRequest(err.Error())
	}
	return nil
}
