package install

import (
	"io"
	"net/http"

	"libs/errorcode"
	"libs/odm"

	"github.com/qiniu/xlog"
)

var (
	bodyLimit = 8 << 20
)

func GetResource(log *xlog.Logger, db *odm.DB, id int64, collName string, req *http.Request) (int, interface{}) {

	v, err := db.Find2(odm.M{"_id": id}, collName)
	if err != nil {
		log.Error(err)
		return errorcode.HandleError(err)
	}

	return http.StatusOK, v
}

func PostResource(log *xlog.Logger, db *odm.DB, collName string, req *http.Request) (int, interface{}) {
	limitReader := io.LimitReader(req.Body, int64(bodyLimit))

	v, err := db.Insert2(collName, limitReader)
	if err != nil {
		log.Error(err)
		return errorcode.HandleError(err)
	}
	return http.StatusCreated, v
}
