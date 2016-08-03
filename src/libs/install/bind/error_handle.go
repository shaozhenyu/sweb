package bind

import (
	"net/http"

	"libs/errorcode"
)

func ErrorHandle(err error, resp http.ResponseWriter) {
	if err != nil {
		status := http.StatusBadRequest
		if e, ok := err.(*errorcode.ErrCodeMessage); ok {
			status = e.HttpStatus
		}
		resp.WriteHeader(status)
		resp.Write([]byte(err.Error()))
	}
}
